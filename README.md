# NewsApp

A Go-based news aggregation application that fetches news articles from NewsAPI and stores them in a SQLite database. The application processes news topics from an input file and generates formatted output files with the fetched articles.

## 🏗️ Architecture

The application uses a Docker-based architecture with dedicated services:

- **sqlite-db**: Creates and manages the SQLite database file
- **newsapp**: Main application that connects to the shared SQLite database
- **sqlite-data**: Persistent Docker volume for database storage

## 🚀 Quick Start

### Prerequisites
- Docker and Docker Compose installed
- NewsAPI key from [https://newsapi.org/](https://newsapi.org/)

### How to Run (Simplest, no build)
If you already have the prebuilt image loaded (via `docker load -i newsapp-image.tar`), run it like this:

```bash
# 1) Prepare a working folder
mkdir -p ~/newsapp-run && cd ~/newsapp-run

# 2) Create input and output
cat > input.txt << 'EOF'
kiwi,2,2
kiwi,2,2
EOF
mkdir -p output

# 3) Set your API key
export NEWS_API_KEY="YOUR_NEWS_API_KEY"

# 4) Create the SQLite volume and DB file (one-time init)
docker run --rm -v newsapp_sqlite-data:/data alpine sh -c "mkdir -p /data && touch /data/newsapp.db"

# 5) Run the app
docker run --rm \
  -e NEWS_API_KEY="$NEWS_API_KEY" \
  -e SQLITE_DB_PATH=/data/newsapp.db \
  -v "$(pwd)/input.txt:/root/input.txt:ro" \
  -v "$(pwd)/output:/root/output" \
  -v newsapp_sqlite-data:/data \
  newsapp-newsapp:latest

# 6) Check results
ls -la output
```

### Setup and Run (with Compose)

1. **Set up environment variables:**
   ```bash
   export NEWS_API_KEY="your_news_api_key_here"
   ```

2. **Run services using the prebuilt image:**
   ```yaml
   # docker-compose.yml (runtime-only)
   services:
     sqlite-db:
       image: alpine:latest
       volumes:
         - sqlite-data:/data
       command: ["sh", "-c", "mkdir -p /data && touch /data/newsapp.db && tail -f /dev/null"]
       restart: no

     newsapp:
       image: newsapp-newsapp:latest
       environment:
         - NEWS_API_KEY=${NEWS_API_KEY}
         - SQLITE_DB_PATH=/data/newsapp.db
       volumes:
         - ./cmd/main/input.txt:/root/input.txt:ro
         - ./cmd/main/output:/root/output
         - sqlite-data:/data
       depends_on:
         - sqlite-db
       restart: no

   volumes:
     sqlite-data:
       driver: local
   ```

   ```bash
   # From the project root (or adapt paths if running elsewhere)
   docker-compose up
   ```

## 📁 Project Structure

```
NewsApp/
├── cmd/
│   └── main/
│       ├── input.txt          # Configuration file with news topics
│       ├── main.go            # Main application entry point
│       └── output/            # Generated news files directory
│           ├── apple.txt
│           └── hello.txt
├── pkg/
│   ├── config/
│   │   └── app.go            # Database configuration
│   ├── controllers/
│   │   └── news-controller.go # Business logic
│   ├── models/
│   │   └── news.go           # Data models and database operations
│   ├── routes/
│   │   └── news-routes.go    # API integration
│   └── utils/
│       └── utils.go          # Utility functions
├── docker-compose.yml        # Docker services configuration
├── Dockerfile               # Docker build instructions
├── go.mod                   # Go module definition
└── go.sum                   # Go dependencies checksums
```

## ⚙️ Configuration

### Environment Variables
- `NEWS_API_KEY`: Required. Your NewsAPI key from https://newsapi.org/
- `SQLITE_DB_PATH`: Database file path (default: `/data/newsapp.db`)

### Input Configuration
Edit `cmd/main/input.txt` to specify news topics:
```
Apple,2,2
hello,2,5
```
Format: `topic,days_back,article_count`

## 🐳 Docker Services

### SQLite Database Service (`sqlite-db`)
- Creates the SQLite database file
- Manages database persistence
- Runs continuously to maintain database availability

### NewsApp Service (`newsapp`)
- Main application
- Connects to shared SQLite database
- Processes news articles and generates output

## 📊 Database Management

### View Database Contents
```bash
# View all tables
docker run --rm -v newsapp_sqlite-data:/data alpine sh -c "apk add --no-cache sqlite && sqlite3 /data/newsapp.db '.tables'"

# View article count
docker run --rm -v newsapp_sqlite-data:/data alpine sh -c "apk add --no-cache sqlite && sqlite3 /data/newsapp.db 'SELECT COUNT(*) FROM articles;'"
```

### Delete Data
```bash
# Delete all articles
docker run --rm -v newsapp_sqlite-data:/data alpine sh -c "apk add --no-cache sqlite && sqlite3 /data/newsapp.db 'DELETE FROM articles;'"

# Delete specific topic
docker run --rm -v newsapp_sqlite-data:/data alpine sh -c "apk add --no-cache sqlite && sqlite3 /data/newsapp.db 'DELETE FROM articles WHERE news_data_id = (SELECT id FROM news_data WHERE searchkey = \"Apple\"); DELETE FROM news_data WHERE searchkey = \"Apple\";'"
```

## 🔧 Manual Docker Commands

### Build and Run
```bash
# Build the image
docker build -t newsapp .

# Run individual services
docker-compose up sqlite-db -d
docker-compose up newsapp

# Run with custom API key
docker run -e NEWS_API_KEY="your_api_key" newsapp-newsapp:latest
```

### Export/Import Images
```bash
# Export images for portability
docker save newsapp-newsapp:latest alpine:latest -o newsapp-complete.tar

# Import on another machine
docker load -i newsapp-complete.tar
```

## 📈 Features

- **Concurrent Processing**: Uses goroutines with controlled parallelism (max 10 concurrent)
- **Database Caching**: Avoids redundant API calls by storing articles in SQLite
- **Duplicate Prevention**: Uses URL-based deduplication
- **Error Handling**: Graceful handling of API failures and database issues
- **Formatted Output**: Generates timestamped files with source tracking
- **Incremental Updates**: Appends new articles to existing records

## 🗄️ Database Schema

- **NewsData table**: Stores search topics
- **Article table**: Stores individual articles with foreign key to NewsData
- **Fields**: Title, Author, URL, Description, PublishedAt

## 📦 Docker Volumes

- `./cmd/main/input.txt:/root/input.txt:ro`: Read-only input configuration
- `./cmd/main/output:/root/output`: Output directory for news files
- `sqlite-data:/data`: Shared volume for SQLite database

## 🎯 Benefits

- **Service Separation**: Database and application are separate services
- **Persistent Data**: Database survives container restarts
- **Scalability**: Can easily add more services that use the same database
- **Portability**: Self-contained Docker images work on any machine
- **No Dependencies**: No need for Go or SQLite on target machine

## 🚀 Performance

- **Image Size**: 55.2MB total (41.9MB app + 13.3MB Alpine base)
- **Fast Startup**: Quick container initialization
- **Low Memory**: Efficient resource usage
- **Concurrent Processing**: Handles multiple topics simultaneously

## 🔍 Troubleshooting

### Common Issues
1. **API Key Missing**: Ensure `NEWS_API_KEY` environment variable is set
2. **Database Connection**: Check if `sqlite-db` service is running
3. **Output Files**: Verify volume mapping in `docker-compose.yml`

### Logs
```bash
# View application logs
docker-compose logs newsapp

# View database service logs
docker-compose logs sqlite-db
```

## 📝 License

This project is open source and available under the MIT License.

# Note:
This project was developed as part of the coursework for MSCS621: Cloud Computing at Marist University, instructed by Rasit Onur Topaloglu, Ph.D. The project specifications and requirements were provided as part of the course assignment.
