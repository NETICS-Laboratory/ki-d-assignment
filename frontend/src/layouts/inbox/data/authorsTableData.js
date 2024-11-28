/* eslint-disable react/prop-types */
// Soft UI Dashboard React components
import SoftBox from "components/SoftBox";
import SoftTypography from "components/SoftTypography";
import SoftAvatar from "components/SoftAvatar";
import SoftBadge from "components/SoftBadge";

// Images
import team2 from "assets/images/team-2.jpg";
import team3 from "assets/images/team-3.jpg";
import team4 from "assets/images/team-4.jpg";

function Author({ image, name, email }) {
  return (
    <SoftBox display="flex" alignItems="center" px={1} py={0.5}>
      <SoftBox mr={2}>
        <SoftAvatar src={image} alt={name} size="sm" variant="rounded" />
      </SoftBox>
      <SoftBox display="flex" flexDirection="column">
        <SoftTypography variant="button" fontWeight="medium">
          {name}
        </SoftTypography>
        <SoftTypography variant="caption" color="secondary">
          {email}
        </SoftTypography>
      </SoftBox>
    </SoftBox>
  );
}

function Function({ job, org }) {
  return (
    <SoftBox display="flex" flexDirection="column">
      <SoftTypography variant="caption" fontWeight="medium" color="text">
        {job}
      </SoftTypography>
      <SoftTypography variant="caption" color="secondary">
        {org}
      </SoftTypography>
    </SoftBox>
  );
}

const authorsTableData = {
  columns: [
    { name: "data", align: "center" },
    { name: "ciphertext", align: "left" },
  ],

  rows: [
    {
      data: (
        <SoftTypography variant="caption" fontWeight="medium">
          Name
        </SoftTypography>
      ),
      ciphertext: (
        <SoftTypography variant="caption" color="secondary" fontWeight="medium">
          Name
        </SoftTypography>
      ),
    },
    {
      data: (
        <SoftTypography variant="button" fontWeight="medium">
          Adress
        </SoftTypography>
      ),
      ciphertext: (
        <SoftTypography variant="caption" color="secondary" fontWeight="medium">
          Name
        </SoftTypography>
      ),
    },
    {
      data: (
        <SoftTypography variant="button" fontWeight="medium">
          IDCARD
        </SoftTypography>
      ),
      ciphertext: (
        <SoftTypography variant="caption" color="secondary" fontWeight="medium">
          Name
        </SoftTypography>
      ),
    },
    {
      data: (
        <SoftTypography variant="button" fontWeight="medium">
          Name
        </SoftTypography>
      ),
      ciphertext: (
        <SoftTypography variant="caption" color="secondary" fontWeight="medium">
          Name
        </SoftTypography>
      ),
    },
    {
      data: (
        <SoftTypography variant="button" fontWeight="medium">
          Name
        </SoftTypography>
      ),
      ciphertext: (
        <SoftTypography variant="caption" color="secondary" fontWeight="medium">
          Name
        </SoftTypography>
      ),
    },
    {
      data: (
        <SoftTypography variant="button" fontWeight="medium">
          Name
        </SoftTypography>
      ),
      ciphertext: (
        <SoftTypography variant="caption" color="secondary" fontWeight="medium">
          Name
        </SoftTypography>
      ),
    },
  ],
};

export default authorsTableData;
