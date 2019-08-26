./build/fleet serve \
    --mysql_address=127.0.0.1:3306 \
    --server_address=0.0.0.0:8080 \
    --mysql_database=kolide \
    --mysql_username=root \
    --mysql_password=123456 \
    --redis_address=127.0.0.1:6379 \
    --server_cert=./cert/server.cert \
    --server_key=./cert/server.key \
    --logging_json \
    --osquery_result_log_file=./osquery_result \
    --osquery_status_log_file=./osquery_status \
    --auth_jwt_key=wmXtTzJXi3V0Z7S9qkDN0C23kJQrGcS/ \
    --event_jpush_id=896e10b745c70105221f2843 \
    --event_jpush_key=b0b0756c29e176c60755362a \
    --server_tls=true


