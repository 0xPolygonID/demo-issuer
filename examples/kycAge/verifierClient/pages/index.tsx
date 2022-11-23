import { useEffect, useState } from "react";
import { QRCode } from "react-qr-svg";
import axios from "axios";
import { Layout } from "../components";
import {Box, Flex, Heading, Paragraph} from "theme-ui";

const Page = (props: {issuerPublicUrl: string, issuerLocalUrl: string}) => {
  const [loading, setLoading] = useState(true);
  const [qrDataSig, setQRDataSig] = useState({});
  const [qrDataMTP, setQRDataMTP] = useState({});
  const [dateData, setDateData] = useState({});

  const checkVerificationStatus = async (sessionID: string) => {
    try {
      const resp = await axios.get(
          "http://" +props.issuerLocalUrl + `/api/v1/status?id=${sessionID}`
      );
      if (resp.status === 200) {
        return true;
      }

      return false;

    } catch (err) {
      console.log("err->", err);
      return false;
    }
  };

  useEffect(() => {
    (async () => {
      const respSig = await axios.get(
          "http://" +props.issuerLocalUrl + "/api/v1/requests/age-kyc?circuitType=credentialAtomicQuerySig"
      );
      const respMTP = await axios.get(
          "http://" +props.issuerLocalUrl + "/api/v1/requests/age-kyc?circuitType=credentialAtomicQueryMTP"
      )


      setQRDataSig(respSig.data);
      setQRDataMTP(respMTP.data)

      const dateLessThan = `${respSig.data.body.scope[0].rules.query.req.birthday.$lt}`;
      const year = dateLessThan.substring(0, 4);
      const month = dateLessThan.substring(4, 6);
      const day = dateLessThan.substring(6, 8);
      const parsedDate = month + "/" + day + "/" + year;

      setDateData(parsedDate);
      setLoading(false);

      const sessionSigID = respSig.headers["x-id"];
      const sessionMtpID = respMTP.headers["x-id"]

      const intervalSig = setInterval(async () => {
        const isVerified = await checkVerificationStatus(sessionSigID);
        if (isVerified) {
          clearInterval(intervalSig);
          alert("verification succeeded with signature proof ✅");
        }
      }, 2000);

      const intervalMtp = setInterval(async () => {
        const isVerified = await checkVerificationStatus(sessionMtpID);
        if (isVerified) {
          clearInterval(intervalMtp);
          alert("verification succeeded with mtp proof ✅");
        }
      }, 2000);
    })();
  }, []);
  return (
    <Layout>
      {loading ? (
        <h1>Loading</h1>
      ) : (
        <Flex
          sx={{ flex: 1, flexDirection: "column", variant: "layout.allCenter" }}
        >
          <Heading sx={{ textAlign: "center", fontSize: [32], my: [4] }}>
            Verify your claim 👮‍♀️
          </Heading>

          <Flex>
            <Box sx={{ margin: "0 30px 0 0" }}>
              <Heading sx={{ textAlign: "center", fontSize: [24] }}>With signature</Heading>
              <QRCode
                level="Q"
                style={{ width: 356 }}
                value={JSON.stringify(qrDataSig)}
              />
            </Box>

            <Box sx={{ margin: "0 0 0 30px" }}>
              <Heading sx={{ textAlign: "center", fontSize: [24] }}>With MTP</Heading>
              <QRCode
                  level="Q"
                  style={{ width: 356 }}
                  value={JSON.stringify(qrDataMTP)}
              />
            </Box>
          </Flex>

          <Paragraph sx={{ variant: "text.para" }}>
            Scan this to verify you were born before the date {dateData}.
          </Paragraph>
        </Flex>
      )}
    </Layout>
  );
};

export async function getServerSideProps(context) {
  const yaml = require('js-yaml');
  const fs = require('fs');

  let  issuerPublicUrl = "";
  let  issuerLocalUrl = "";


  try {
    const doc = yaml.load(fs.readFileSync('./../../../issuer/issuer_config.default.yaml', 'utf8'));
    issuerPublicUrl = doc.public_url;
    issuerLocalUrl = doc.local_url;

  } catch (e) {
    console.log("encounter error on load config file, err: " + e);
    process.exit(1);
  }

  return {
    props: {issuerPublicUrl: issuerPublicUrl, issuerLocalUrl: issuerLocalUrl },
  }
}

export default Page;
