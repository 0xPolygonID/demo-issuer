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
    "version": Math.floor(Math.random() * 1000) + 1,
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