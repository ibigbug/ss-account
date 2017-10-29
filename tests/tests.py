from BaseHTTPServer import BaseHTTPRequestHandler, HTTPServer


class LoadproviderHandler(BaseHTTPRequestHandler):
    def do_GET(self):
        self.send_response(200)
        self.send_header('Content-Type', 'text/plain')
        self.end_headers()

        # 512 Mb
        for i in range(1024 * 512):
            self.wfile.write('x' * 1024)


def run(server_class=HTTPServer, handler_class=LoadproviderHandler, port=8080):
    server_address = ('', port)
    httpd = server_class(server_address, handler_class)
    print('Starting httpd at %s:%s' % server_address)
    httpd.serve_forever()


if __name__ == '__main__':
    run()
