# Leaderboard API Documentation

## Overview
The Leaderboard API retrieves player rankings based on their total points earned from exploring treasure chests within a specific source/category.

## Endpoint
```
POST /leaderboard
```

## Request Body

| Field    | Type   | Required | Description                                    | Constraints           |
|----------|--------|----------|------------------------------------------------|-----------------------|
| `source` | string | Yes      | Source/category of the treasure hunt           | Non-empty             |
| `limit`  | number | Yes      | Number of leaderboard entries to return        | Min: 1, Max: 100      |
| `offset` | number | No       | Offset for pagination (default: 0)             | Min: 0                |

### Example Request
```json
{
  "source": "test",
  "limit": 10,
  "offset": 0
}
```

## Response Body

### Success Response (200 OK)

| Field         | Type                   | Description                                    |
|---------------|------------------------|------------------------------------------------|
| `leaderboard` | array[LeaderboardEntry]| List of leaderboard entries                    |
| `total`       | number                 | Total number of players in this source         |
| `source`      | string                 | Source/category of the treasure hunt           |

#### LeaderboardEntry Object

| Field    | Type   | Description                                |
|----------|--------|--------------------------------------------|
| `name`   | string | Player's name                              |
| `rank`   | number | Player's rank (1 = highest points)         |
| `points` | number | Total points earned by the player          |

### Example Success Response
```json
{
  "leaderboard": [
    {
      "name": "Alice Johnson",
      "rank": 1,
      "points": 500
    },
    {
      "name": "Bob Smith",
      "rank": 2,
      "points": 450
    },
    {
      "name": "Charlie Brown",
      "rank": 3,
      "points": 400
    }
  ],
  "total": 25,
  "source": "test"
}
```

### Error Response (400 Bad Request)
```json
{
  "success": false,
  "error": "Invalid request body"
}
```

### Error Response (500 Internal Server Error)
```json
{
  "success": false,
  "error": "Failed to get leaderboard: database error"
}
```

## How It Works

1. **Query Explorer Data**: The API queries the `treasure_explorer` table for all entries matching the specified source
2. **Aggregate Scores**: Groups results by player and sums their scores
3. **Calculate Ranks**: Assigns ranks based on total points (higher points = lower rank number)
4. **Apply Pagination**: Returns only the requested page of results (limit + offset)
5. **Return Results**: Includes leaderboard entries, total player count, and source

## Performance Considerations

### Database Query Optimization

The leaderboard query performs the following operations:
```sql
SELECT 
  p.name,
  SUM(te.score) as total_points,
  RANK() OVER (ORDER BY SUM(te.score) DESC) as rank
FROM treasure_explorer te
JOIN players p ON te.player_id = p.id
WHERE te.source = ?
GROUP BY te.player_id, p.name
ORDER BY total_points DESC
LIMIT ? OFFSET ?
```

### Is it OK to Run GROUP BY Frequently?

**Short Answer**: Yes, with proper indexing, it's generally fine for most use cases.

**Considerations**:

1. **Index on `source` column**: ✅ **Added** 
   - The `idx_treasure_explorer_source` index filters rows efficiently
   - This is the most important index for this query

2. **Index on `player_id` column**: ✅ **Already exists**
   - The `idx_treasure_explorer_player_id` index helps with joins

3. **Query Performance Factors**:
   - **Number of explorers per source**: 
     - < 1,000: Excellent performance (< 10ms)
     - 1,000 - 10,000: Good performance (10-50ms)
     - 10,000 - 100,000: Acceptable (50-200ms)
     - > 100,000: May need optimization
   
4. **Optimization Strategies** (if needed):

   **A. Materialized Views** (Best for frequently-accessed leaderboards):
   ```sql
   CREATE MATERIALIZED VIEW leaderboard_cache AS
   SELECT 
     te.source,
     p.name,
     te.player_id,
     SUM(te.score) as total_points,
     RANK() OVER (PARTITION BY te.source ORDER BY SUM(te.score) DESC) as rank
   FROM treasure_explorer te
   JOIN players p ON te.player_id = p.id
   GROUP BY te.source, te.player_id, p.name;
   
   -- Refresh strategy: every 1-5 minutes or after chest claims
   REFRESH MATERIALIZED VIEW CONCURRENTLY leaderboard_cache;
   ```

   **B. Caching Layer** (Redis/Memcached):
   - Cache leaderboard results for 30-60 seconds
   - Invalidate cache when new chests are claimed
   - Reduces database load significantly

   **C. Background Jobs**:
   - Pre-calculate leaderboards every N minutes
   - Store results in a separate `leaderboard_snapshot` table
   - Serve from snapshot table instead of live calculation

   **D. Composite Index** (if JOIN is slow):
   ```sql
   CREATE INDEX idx_explorer_source_player ON treasure_explorer(source, player_id);
   ```

### Recommended Approach by Scale

| Daily Active Users | Approach                                      | Expected Latency |
|--------------------|-----------------------------------------------|------------------|
| < 100              | Direct query (current implementation)         | < 10ms           |
| 100 - 1,000        | Direct query with monitoring                  | < 50ms           |
| 1,000 - 10,000     | Add Redis cache (60s TTL)                    | < 5ms (cached)   |
| 10,000+            | Materialized view + Redis                     | < 2ms (cached)   |

### Current Implementation Status

✅ **Optimized for most use cases** (< 10,000 active players per source)
- Indexed `source` column
- Indexed `player_id` column
- Efficient SQL with early filtering

🔧 **Ready for scale** - If you experience slow queries:
1. Add Redis caching first (easiest)
2. Create materialized views if still slow
3. Monitor query performance with `EXPLAIN ANALYZE`

## Testing

Use the provided test script:
```bash
./api-curls/04-leaderboard.sh
```

Or test manually:
```bash
curl -X POST 'http://localhost:8080/leaderboard' \
  -H 'Content-Type: application/json' \
  -d '{
    "source": "test",
    "limit": 10,
    "offset": 0
  }'
```

## Example Use Cases

### 1. Display Top 10 Players
```json
{
  "source": "event-2024",
  "limit": 10,
  "offset": 0
}
```

### 2. Paginated Leaderboard (Next Page)
```json
{
  "source": "event-2024",
  "limit": 10,
  "offset": 10
}
```

### 3. Full Leaderboard Export
```json
{
  "source": "event-2024",
  "limit": 100,
  "offset": 0
}
```

## Notes

- Ranks are calculated dynamically based on total points
- Players with the same points receive the same rank
- The `total` field shows total unique players, not total entries
- Empty sources return empty leaderboard array with total = 0
- Limit is automatically capped at 100 to prevent excessive data transfer

