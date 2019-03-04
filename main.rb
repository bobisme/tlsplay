#!/usr/bin/env ruby
require 'openssl'
require 'net/http'
require 'json'

class Client
  DEFAULT_OPTIONS = {
    use_ssl: true,
    # TRUSTING THE SERVER
    verify_mode: OpenSSL::SSL::VERIFY_NONE,
    keep_alive_timeout: 30,
    cert: OpenSSL::X509::Certificate.new(File.read('./certs/service-1234@accounts.example.com.pem')),
    key: OpenSSL::PKey::EC.new(File.read('./certs/service-1234@accounts.example.com-key.pem'))
  }

  def initialize
    @http = Net::HTTP.start("localhost", 8443, DEFAULT_OPTIONS)
  end

  def fetch
    response = @http.request Net::HTTP::Get.new "/echo"
    response.body
  end
end

puts "RUBY GOT: #{Client.new.fetch}"
