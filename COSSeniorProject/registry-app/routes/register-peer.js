var express = require('express');
var router = express.Router();
var Response = require("../types/response/response");
var redisService =require("../services /redis/redis");
var middleWare = require("../midleware/enable-cors");


router.post("/register", middleWare.enableCors(req,res), (req, res)=>{

    const peerAddress = req.body.peerID;
    const redisClient = redisService.createRedisClient();
    redisService.put(redisClient,peerAddress);
    const resp = new Response();
    return res.send(
        resp.getOpStatus(
            "success",
            false,
            "operation was carried out successfully"
        )
    );
});