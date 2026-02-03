"""Sample Flask application with PostgreSQL and Redis for discovery testing."""
import os
from flask import Flask, jsonify
import psycopg2
import redis

app = Flask(__name__)

# Database configuration
DATABASE_URL = os.getenv('DATABASE_URL', 'postgresql://localhost:5432/mydb')
REDIS_URL = os.getenv('REDIS_URL', 'redis://localhost:6379')

# Initialize Redis client
redis_client = redis.from_url(REDIS_URL)


def get_db_connection():
    """Create a database connection."""
    return psycopg2.connect(DATABASE_URL)


@app.route('/')
def index():
    """Health check endpoint."""
    return jsonify({'status': 'ok'})


@app.route('/api/data')
def get_data():
    """Get data from database with Redis caching."""
    cache_key = 'api:data'
    
    # Try cache first
    cached = redis_client.get(cache_key)
    if cached:
        return jsonify({'data': cached.decode(), 'source': 'cache'})
    
    # Query database
    conn = get_db_connection()
    cur = conn.cursor()
    cur.execute('SELECT * FROM data LIMIT 10')
    rows = cur.fetchall()
    cur.close()
    conn.close()
    
    # Cache result
    result = str(rows)
    redis_client.setex(cache_key, 300, result)
    
    return jsonify({'data': result, 'source': 'database'})


if __name__ == '__main__':
    port = int(os.getenv('PORT', 5000))
    app.run(host='0.0.0.0', port=port)
