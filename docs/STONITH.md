# STONITH Mechanism in Kubegres

## What is STONITH?

STONITH (Shoot The Other Node In The Head) is a fencing mechanism used in high-availability clusters to ensure that a failed or unresponsive node cannot cause data corruption or inconsistencies. In PostgreSQL context, it ensures that only one primary instance is active at any time, preventing split-brain scenarios.

## Split-Brain in PostgreSQL

A split-brain scenario occurs when two instances believe they are the primary database server. This can happen during network partitions or when the original primary becomes unreachable but is still running. In such situations, both primaries might accept writes, leading to data divergence and potential data loss when the cluster heals.

## STONITH in Kubegres

Starting from version X.Y.Z, Kubegres implements a STONITH mechanism to ensure safe failover of PostgreSQL primary instances. When enabled, Kubegres explicitly terminates the old primary pod before promoting a replica to become the new primary.

### How It Works

1. When Kubegres detects that a primary pod is unhealthy or unavailable, it initiates the failover process.
2. If STONITH is enabled, Kubegres forcefully deletes the old primary pod with a zero grace period.
3. Kubegres waits for confirmation that the old primary is fully terminated (up to 30 seconds).
4. Only after the old primary is terminated, Kubegres promotes a replica to become the new primary.
5. A new replica pod is created to replace the terminated primary.

### Configuration

To enable STONITH in your Kubegres deployment, add `enableSTONITH: true` to the failover section of your Kubegres spec:
