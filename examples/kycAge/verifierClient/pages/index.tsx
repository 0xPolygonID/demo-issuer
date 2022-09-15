import { useEffect, useState } from "react";
import { QRCode } from "react-qr-svg";
import axios from "axios";
import { makeClaimRequest } from "../utils/utils";
import { useRouter } from "next/router";
import { Layout } from "../components";
import { Flex, Heading, Paragraph } from "theme-ui";

const Page = () => {
  const [loading, setLoading] = useState(true);
  const [qrData, setQRData] = useState({});

  const router = useRouter();

  const checkAuthStatus = async (sessionID: string) => {
    try {
      const resp = await axios.get(
        `http://localhost:3000/api/status?id=${sessionID}`
      );

      const userID = resp.data.id;

      if (userID) {
        const resp = await axios(makeClaimRequest(userID));
        // TODO: Error Handling
        const claimID = resp.data.id ? resp.data.id : "";
        return { claimID, userID };
      }
    } catch (err) {
      // TODO: Error Handling
      console.log("err->", err);
      return false;
    }
  };

  useEffect(() => {
    (async () => {
      const resp = await axios.get(
        "http://localhost:3000/api/sign-in?type=random"
      );

      setQRData(resp.data);
      setLoading(false);

      const sessionID = resp.headers["x-id"];

      const interval = setInterval(async () => {
        const resp = await checkAuthStatus(sessionID);
        if (resp) {
          clearInterval(interval);
          alert("verification succeded ✅");
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

          <QRCode
            level="Q"
            style={{ width: 256 }}
            value={JSON.stringify(qrData)}
          />

          <Paragraph sx={{ variant: "text.para" }}>
            Scan this to verify you are above 22 years old.
          </Paragraph>
        </Flex>
      )}
    </Layout>
  );
};

export default Page;
