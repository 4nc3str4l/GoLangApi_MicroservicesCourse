#!/bin/bash
protoc ./blog_pb/blog.proto --go_out=plugins=grpc:.