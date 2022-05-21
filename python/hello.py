#! /usr/bin/env python3
# -*- coding: utf-8 -*-
import grpc
import time
from concurrent import futures
import os, base64
import sys
import search_pb2 as proto
import search_pb2_grpc as pb2
import ddddocr
import urllib.request

_ONE_DAY_IN_SECONDS = 60 * 60 * 24
_HOST = '0.0.0.0'
_PORT = '7000'

class SearchService(pb2.SearchServiceServicer):
    def Search(self, request, context):
        ocr = ddddocr.DdddOcr(old=True)
        str = request.request
        data = str.split(',')[1]
        image_data = base64.b64decode(data)
        res = ocr.classification(image_data)
        print(res)
        return proto.SearchRequest(request=res)
def serve():
    grpcServer = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    pb2.add_SearchServiceServicer_to_server(SearchService(), grpcServer)
    grpcServer.add_insecure_port(_HOST + ':' + _PORT)
    grpcServer.start()
    print("***   "+_HOST+':'+_PORT+"端口开始占用！！       **\n" +
          "*       验证码破解启动成功！！            *\n" +
          "*       永不宕机     永无BUG              *\n" +
          "*     本系统由gva团队精心打造             *\n" +
          "*******************************************\n"
    )
    try:
        while True:
            time.sleep(_ONE_DAY_IN_SECONDS)
    except KeyboardInterrupt:
        grpcServer.stop(0)
if __name__ == '__main__':
    serve()