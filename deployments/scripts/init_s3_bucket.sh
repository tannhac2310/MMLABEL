#!/bin/sh
mc config host add minio http://s3_storage:9000 ${MINIO_ACCESS_KEY} ${MINIO_SECRET_KEY}
mc mb minio/${INIT_BUCKET_NAME}
mc policy download minio/${INIT_BUCKET_NAME}
