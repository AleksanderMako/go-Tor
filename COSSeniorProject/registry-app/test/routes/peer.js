var redis = require("../../services /redis/redis");
var expect = require("chai").expect;
var mocha = require("mocha");

const supertest = require("supertest");
const api = supertest("http://tor-registry:4500");

describe("peer routes ", ()=>{

    it("should return 200 when all data is ok for peer registration ", async ()=>{
        //Arrange 

        const peerRegistrationRequest = {
            peerID:"peer1:9000"
        }
        const registerPeerRoute = "/api/register/"
        //Act 

        const res = await api.post("/peer")
            .set("Accept", "application/json")
            .send(peerRegistrationRequest);

        //Assert 
         
        expect(res.status).to.eql(200);



    });
});