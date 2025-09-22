# NewsApp Docker Setup with SQLite Service

## Architecture
The application now uses a dedicated SQLite service with proper volume management:

- **sqlite-db**: Creates and manages the SQLite database file
- **sqlite-web**: Provides a web interface for database management (port 8080)
- **newsapp**: Main application that connects to the shared SQLite database

## Prerequisites
- Docker and Docker Compose installed
- NewsAPI key from https://newsapi.org/

## Quick Start

1. **Set up environment variables:**
   ```bash
   export NEWS_API_KEY="your_news_api_key_here"
   ```

2. **Build and run all services:**
   ```bash
   docker-compose up --build
   ```

3. **Run in detached mode:**
   ```bash
   docker-compose up -d --build
   ```

## Services

### SQLite Database Service (`sqlite-db`)
- Creates the SQLite database file
- Manages database persistence
- Runs continuously to maintain database availability

### SQLite Web Interface (`sqlite-web`)
- Provides web-based database management
- Accessible at: http://localhost:8080
- View tables, run queries, manage data

### NewsApp Service (`newsapp`)
- Main application
- Connects to shared SQLite database
- Processes news articles and generates output

## Manual Docker Commands

1. **Build the image:**
   ```bash
   docker build -t newsapp .
   ```

2. **Run individual services:**
   ```bash
   # Start SQLite database
   docker-compose up sqlite-db -d
   
   # Start web interface
   docker-compose up sqlite-web -d
   
   # Start main application
   docker-compose up newsapp
   ```

## File Structure
- `cmd/main/input.txt`: Configuration file with news topics
- `cmd/main/output/`: Directory where generated news files are stored
- `data/`: Directory for SQLite database file (managed by Docker volume)

## Environment Variables
- `NEWS_API_KEY`: Required. Your NewsAPI key from https://newsapi.org/
- `SQLITE_DB_PATH`: Database file path (default: `/data/newsapp.db`)

## Volumes
- `./cmd/main/input.txt:/root/input.txt:ro`: Read-only input configuration
- `./cmd/main/output:/root/output`: Output directory for news files
- `sqlite-data:/data`: Shared volume for SQLite database

## Database Management
- **Web Interface**: Access http://localhost:8080 to manage the database
- **Persistent Storage**: Database persists in Docker volume `sqlite-data`
- **Backup**: Volume can be backed up using Docker volume commands

## Benefits
- **Service Separation**: Database and application are separate services
- **Web Management**: Easy database management through web interface
- **Persistent Data**: Database survives container restarts
- **Scalability**: Can easily add more services that use the same database
- **Monitoring**: Web interface allows real-time database monitoring
