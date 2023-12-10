#!/bin/sh
mc config host add minio http://s3_storage:9000 ${MINIO_ACCESS_KEY} ${MINIO_SECRET_KEY}
mc mb minio/${INIT_BUCKET_NAME}
mc anonymous download minio/mm-printing-images
mc anonymous set download mm-printing-images
mc anonymous set public mm-printing-images
mc anonymous list mm-printing-images

mc policy download minio/${INIT_BUCKET_NAME}
  mc anonymous list mm-printing-images