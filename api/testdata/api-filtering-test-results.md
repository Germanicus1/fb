# API Filtering Capabilities Test Results

**Test Time:** 2026-02-13T23:05:12+01:00

**Baseline Ticket Count:** 64
**Baseline Response Size:** 37074 bytes

## Conclusion

Server-side filtering is NOT supported. All tested parameters either return errors or do not affect results. Client-side filtering is required.

## Parameter Test Results

| Parameter | Value | Accepted | Ticket Count | Response Size | Affects Results | Error |
|-----------|-------|----------|--------------|---------------|-----------------|-------|
| board | test-board | Yes | 64 | 37074 | No | - |
| bin | test-bin | Yes | 64 | 37074 | No | - |
| boardId | test-board-id | Yes | 64 | 37074 | No | - |
| binId | test-bin-id | Yes | 64 | 37074 | No | - |
| board_id | test-board-id | Yes | 64 | 37074 | No | - |
| bin_id | test-bin-id | Yes | 64 | 37074 | No | - |
| boardName | test-board-name | Yes | 64 | 37074 | No | - |
| binName | test-bin-name | Yes | 64 | 37074 | No | - |
