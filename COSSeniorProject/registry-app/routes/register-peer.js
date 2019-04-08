var express = require('express');
var router = express.Router();
var Response = require("../types/response/response");
var redisService = require("../services /redis/redis");
var middleWare = require("../midleware/enable-cors");


//router.use(middleWare.enableCors);
router.post("/", (req, res) => {
    // peers register under the peer set 
    const peerAddress = req.body.peerID;
    if (!peerAddress) {
        return res.status(500).send("Missing peer id ")
    }
    const redisClient = redisService.createRedisClient();
    const setKey = "peers"
    redisService.putInSet(redisClient, setKey, peerAddress, "available");
    console.log(peerAddress);
    const resp = Response.getOpStatus("success", false, "operation was carried out successfully")
    return res.send(JSON.stringify(resp));
});

router.get("/:setkey", async (req, res) => {

    const setKey = req.params.setkey;
    if (!setKey) {
        return res.status(500).send("MISSING SET KEY ")
    }
    let redisClient = redisService.createRedisClient();
    const setItems = await redisService.getSetKeys(redisClient, setKey);
    let freepeers = await redisService.getFreePeers(setItems)
    const response = {
        peers: freepeers
    }
    return res.json(response);
});
router.post("/test", (req, res) => {

    console.log(req.body);
    return res.send("message received");
});
module.exports = router;