from http.server import BaseHTTPRequestHandler, HTTPServer

class AlertListenerHandler(BaseHTTPRequestHandler):
    def do_POST(self):
        # 1. Verify the request is targeting the correct endpoint
        if self.path == '/alert':
            # 2. Extract the Content-Length header to know how many bytes to read
            content_length = int(self.headers.get('Content-Length', 0))

            # 3. Read the exact number of raw bytes from the stream
            raw_body = self.rfile.read(content_length)

            # 4. Decode the bytes into a standard UTF-8 string
            decoded_alert = raw_body.decode('utf-8')

            # 5. Print the output to your terminal
            print("\n🚨 [ALERT RECEIVED] 🚨")
            print(decoded_alert)
            print("-" * 30)

            # 6. Send a standard HTTP 200 OK response back to your Go agent
            self.send_response(200)
            self.send_header('Content-Type', 'text/plain')
            self.end_headers()
            self.wfile.write(b"Alert logged successfully")
        else:
            # Handle non-existent paths gracefully
            self.send_response(404)
            self.end_headers()

def run_server():
    server_address = ('localhost', 4444)
    httpd = HTTPServer(server_address, AlertListenerHandler)
    print("🚀 Python listener is up and running on http://localhost:4444...")
    try:
        httpd.serve_forever()
    except KeyboardInterrupt:
        print("\nStopping Python listener.")
        httpd.server_close()

if __name__ == '__main__':
    run_server()
