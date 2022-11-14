import { issuerID } from "./constants";
import axios from 'axios'

// export const checkAuthStatus =  async (sessionID:string, props: {issuerPublicUrl: string, issuerLocalUrl: string}) => {
//      try {
//      const resp = await axios.get("http://" + props.issuerLocalUrl + `/api/v1/status?id=${sessionID}`)
//
//      console.log('Here', resp.data)
//
//      const userID = resp.data.id
//
//      if(userID){
//
//         const resp = await axios(makeAgeClaimRequest(userID, props))
//         // TODO: Error Handling
//         const claimID = resp.data.id ? resp.data.id : ""
//         return {claimID, userID}
//      }
//     }
//     // TODO: Error Handling
//     catch (err){
//       console.log('err->', err)
//     return false
//     }
// }



// export const makeAgeClaimData = (claimID:string, userID:string, props: {issuerPublicUrl: string, issuerLocalUrl: string}) => {
//   return{
//       id:"f7a3fae9-ecf1-4603-804e-8ff1c7632636",
//       typ:"application/iden3comm-plain-json",
//       type:"https://iden3-communication.io/credentials/1.0/offer",
//       thid:"f7a3fae9-ecf1-4603-804e-8ff1c7632636",
//       body:{url: props.issuerPublicUrl + `/api/v1/agent`,
//       credentials:[{"id":claimID,"description":"KYCAgeCredential"}]},
//       from: userID,
//       to: issuerID
//   }
// }



export const makeAgeClaimRequest = (dob: number, userID:string, props: {issuerPublicUrl: string, issuerLocalUrl: string}) => {
    const data = JSON.stringify({
    "identifier": userID,
    "schema": {
      "url": "https://raw.githubusercontent.com/iden3/claim-schema-vocab/main/schemas/json-ld/kyc-v2.json-ld",
      "type": "KYCAgeCredential"
    },
    "data": {
      "birthday": dob, //19960424,
      "documentType": 1
    },
    "expiration": 12345678888
  });
  
  const config = {
    method: 'post',
    url: "http://" + props.issuerLocalUrl + '/api/v1/claims',
    headers: { 
      'Content-Type': 'application/json'
    },
    data : data
  }

    return config
};