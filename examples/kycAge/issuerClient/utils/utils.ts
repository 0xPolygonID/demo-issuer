import { agentEndPoint, issuerID } from "./constants";
import axios from 'axios'

export const checkAuthStatus =  async (sessionID:string) => {
     try {
     const resp = await axios.get(`http://localhost:3000/api/status?id=${sessionID}`) 

     console.log('Here', resp.data)

     const userID = resp.data.id

     if(userID){
         
        const resp = await axios(makeClaimRequest(userID))
        // TODO: Error Handling
        const claimID = resp.data.id ? resp.data.id : ""
        return {claimID, userID}
     }
    }
    // TODO: Error Handling
    catch (err){
      console.log('err->', err)
    return false
    }
  }

export const makeAgeClaimData = (claimID:string, userID:string) => {
  return{
  id:"f7a3fae9-ecf1-4603-804e-8ff1c7632636",
  typ:"application/iden3comm-plain-json",
  type:"https://iden3-communication.io/credentials/1.0/offer",
  thid:"f7a3fae9-ecf1-4603-804e-8ff1c7632636",
  body:{url:agentEndPoint, 
  credentials:[{"id":claimID,"description":"KYCAgeCredential"}]},
  from: userID,
  to: issuerID
}}

export const makeClaimRequest = (userID:string) => { 
    
    const data = JSON.stringify({
    "identifier": userID,
    "schema": {
      "url": "https://raw.githubusercontent.com/iden3/claim-schema-vocab/main/schemas/json-ld/kyc-v2.json-ld",
      "type": "KYCAgeCredential"
    },
    "data": {
      "birthday": 19960424,
      "documentType": 1
    },
    "expiration": 12345678888
  });
  
  const config = {
    method: 'post',
    url: 'http://localhost:8001/api/v1/claims',
    headers: { 
      'Content-Type': 'application/json'
    },
    data : data
  }

    return config
};