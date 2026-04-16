import subprocess
import hashlib
import base64
import os

# Generate RSA private key
subprocess.run(['openssl', 'genrsa', '-out', 'ext.pem', '2048'], stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)

# Get public key in DER format
result = subprocess.run(['openssl', 'rsa', '-in', 'ext.pem', '-pubout', '-outform', 'DER'], capture_output=True)
der_pub = result.stdout

# Base64 encode for manifest.json
b64_pub = base64.b64encode(der_pub).decode('utf-8')

# Calculate Extension ID
# SHA256 of DER public key
m = hashlib.sha256()
m.update(der_pub)
hash_hex = m.hexdigest()

# Translate first 32 chars from hex to a-p
trans = str.maketrans('0123456789abcdef', 'abcdefghijklmnop')
ext_id = hash_hex[:32].translate(trans)

print(f"EXTENSION_ID: {ext_id}")
print(f"MANIFEST_KEY: {b64_pub}")
