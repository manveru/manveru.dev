#!/usr/bin/env ruby

require 'ramaze'

# FCGI doesn't like you writing to stdout
Ramaze::Log.loggers = [Ramaze::Logger::Informer.new(__DIR__("../ramaze.fcgi.log"))]

require_relative '../app'

Ramaze.options.trap = :SIGTERM
Ramaze.start(root: __DIR__('../'), started: true)

socket = File.join(*ENV.values_at('RAMAZE_SOCKET', 'RAMAZE_SOCKET_NUMBER'))
socket << '.sock'
puts "Connecting to #{socket}"
FileUtils.rm_f(socket)

Thread.new do
  begin
    sleep 0.5 until File.socket?(socket)
    File.chmod(0660, socket)
    File.chown(Etc.getpwnam('manveru').uid, Etc.getgrnam('www-data').gid, socket)
  rescue Exception => ex
    Ramaze::Log.error(ex)
    raise(ex)
  end
end

Rack::Handler.get(:fastcgi).run(Ramaze, File: socket)
