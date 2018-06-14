package controllers

import (
    "github.com/ninjadotorg/handshake-dispatcher/services"
)

var solrService = new(services.SolrService)
var uploadService = new(services.UploadService)
var fcmService = new(services.FCMService)
var ethereumService = new(services.EthereumService)
