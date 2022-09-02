🎣 catfish
-----
Useful dummy server used for development.

# docker

```bash
docker build . -t catfish
docker run -p 8080:8080 -v ${YOUR_CONFIG}:/etc/catfish/config.yml catfish
```
