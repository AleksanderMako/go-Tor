

var redis = require("../../services /redis/redis");
var chai = require ("chai");
var mocha = require("mocha");
describe("redis service initial tests !", ()=>{

    it ("should put data correctly ", async()=>{

        // Arrange 
        const test = {
            redisData: "this_data"
        }
        const client = redis.createRedisClient();

        //Act 
        redis.put(client, test);

        //Assert 

        const data = await redis.get(client, "testKey");
        chai.expect(data.redisData).equal(test.redisData);
        
    });
   

})