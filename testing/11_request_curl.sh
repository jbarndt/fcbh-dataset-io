
curl -X POST http://localhost:8080 \
-F "audio=@audio.mp3;type=audio/mpeg" \
-F "yaml=@request.yaml;type=application/x-yaml" \
-H "Accept: application/json"
