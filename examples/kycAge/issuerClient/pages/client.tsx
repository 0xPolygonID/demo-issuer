import { Box, Flex, Heading, Paragraph } from "theme-ui";
import { useRouter } from "next/router";
import { makeAgeClaimData } from "../utils/utils";
import { Layout, QRCode } from "../components";

const Page = () => {
  const router = useRouter();
  const claimID = router.query.claimID;
  const userID = router.query.userID;

  let qrData;

  if (typeof claimID === "string" && typeof userID === "string") {
    qrData = makeAgeClaimData(claimID, userID);
  }

  return (
    <Layout>
      <Flex
        sx={{
          flex: 1,
          flexDirection: "column",
          justifyContent: "center",
          alignItems: "center",
        }}
      >
        <Heading sx={{ textAlign: "center", fontSize: [32], my: [4] }}>
          Get Your Claim 🚀
        </Heading>
        <QRCode
          level="Q"
          style={{ width: "100%" }}
          value={JSON.stringify(qrData)}
        />
        <Box>
          <Paragraph sx={{ variant: "text.para" }}>
            This claim proves you are born on January 1st, 2002
          </Paragraph>
        </Box>
      </Flex>
    </Layout>
  );
};

export default Page;
