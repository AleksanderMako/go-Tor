var express = require('express');
var router = express.Router();
var Response = require("../types/response/response");
var redisService = require("../services /redis/redis");
var middleWare = require("../midleware/enable-cors");


//router.use(middleWare.enableCors);
router.post("/", (req, res) => {

    const peerAddress = req.body.peerID;
    const redisClient = redisService.createRedisClient();
    redisService.put(redisClient, peerAddress);
    const resp = Response.getOpStatus("success", false, "operation was carried out successfully")
    return res.send(JSON.stringify(resp));
});

module.exports = router;