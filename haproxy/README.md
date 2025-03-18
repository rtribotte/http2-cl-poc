HAProxy Configuration
Directory Structure

To organize the files, create the following structure:

haproxy/
├── haproxy.cfg
├── certs/
│   ├── server.pem

Combined Certificate and Key

HAProxy requires a combined file containing both the certificate and the private key:

cat server-key.pem server.pem > certs/server.pem

