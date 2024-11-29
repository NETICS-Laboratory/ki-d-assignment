import React, { useEffect, useState } from "react";
import Card from "@mui/material/Card";
import SoftBox from "components/SoftBox";
import SoftTypography from "components/SoftTypography";
import DashboardLayout from "examples/LayoutContainers/DashboardLayout";
import DashboardNavbar from "examples/Navbars/DashboardNavbar";
import Footer from "examples/Footer";
import Table from "examples/Tables/Table";
import axios from "axios";

function Encrypted() {
  const [data, setData] = useState({});
  const [aesRows, setAesRows] = useState([]);
  const [rc4Rows, setRc4Rows] = useState([]);
  const [desRows, setDesRows] = useState([]);

  // Fetch data function
  const fetchData = async () => {
    try {
      const token = localStorage.getItem("token");

      const response = await axios.get("http://127.0.0.1:8090/api/user/me", {
        headers: {
          Authorization: `Bearer ${token}`, // Add the Bearer token here
        },
      });

      const responseData = response.data.data;

      setData(responseData);

      // Set AES rows
      setAesRows([
        {
          data: (
            <SoftTypography variant="caption" fontWeight="medium">
              Name
            </SoftTypography>
          ),
          ciphertext: (
            <SoftTypography variant="caption" color="secondary" fontWeight="medium">
              {responseData.name_aes}
            </SoftTypography>
          ),
        },
        {
          data: (
            <SoftTypography variant="caption" fontWeight="medium">
              Email
            </SoftTypography>
          ),
          ciphertext: (
            <SoftTypography variant="caption" color="secondary" fontWeight="medium">
              {responseData.email_aes}
            </SoftTypography>
          ),
        },
        {
          data: (
            <SoftTypography variant="caption" fontWeight="medium">
              Phone
            </SoftTypography>
          ),
          ciphertext: (
            <SoftTypography variant="caption" color="secondary" fontWeight="medium">
              {responseData.notelp_aes}
            </SoftTypography>
          ),
        },
        {
          data: (
            <SoftTypography variant="caption" fontWeight="medium">
              Address
            </SoftTypography>
          ),
          ciphertext: (
            <SoftTypography variant="caption" color="secondary" fontWeight="medium">
              {responseData.address_aes}
            </SoftTypography>
          ),
        },
        {
          data: (
            <SoftTypography variant="caption" fontWeight="medium">
              ID Card
            </SoftTypography>
          ),
          ciphertext: (
            <SoftTypography variant="caption" color="secondary" fontWeight="medium">
              {responseData.id_card_aes}
            </SoftTypography>
          ),
        },
        {
          data: (
            <SoftTypography variant="caption" fontWeight="medium">
              Username
            </SoftTypography>
          ),
          ciphertext: (
            <SoftTypography variant="caption" color="secondary" fontWeight="medium">
              {responseData.username_aes}
            </SoftTypography>
          ),
        },
        {
          data: (
            <SoftTypography variant="caption" fontWeight="medium">
              Password
            </SoftTypography>
          ),
          ciphertext: (
            <SoftTypography variant="caption" color="secondary" fontWeight="medium">
              {responseData.password_aes}
            </SoftTypography>
          ),
        },
      ]);

      // Set RC4 rows
      setRc4Rows([
        {
          data: (
            <SoftTypography variant="caption" fontWeight="medium">
              Name
            </SoftTypography>
          ),
          ciphertext: (
            <SoftTypography variant="caption" color="secondary" fontWeight="medium">
              {responseData.name_rc4}
            </SoftTypography>
          ),
        },
        {
          data: (
            <SoftTypography variant="caption" fontWeight="medium">
              Email
            </SoftTypography>
          ),
          ciphertext: (
            <SoftTypography variant="caption" color="secondary" fontWeight="medium">
              {responseData.email_rc4}
            </SoftTypography>
          ),
        },
        {
          data: (
            <SoftTypography variant="caption" fontWeight="medium">
              Phone
            </SoftTypography>
          ),
          ciphertext: (
            <SoftTypography variant="caption" color="secondary" fontWeight="medium">
              {responseData.notelp_rc4}
            </SoftTypography>
          ),
        },
        {
          data: (
            <SoftTypography variant="caption" fontWeight="medium">
              Address
            </SoftTypography>
          ),
          ciphertext: (
            <SoftTypography variant="caption" color="secondary" fontWeight="medium">
              {responseData.address_rc4}
            </SoftTypography>
          ),
        },
        {
          data: (
            <SoftTypography variant="caption" fontWeight="medium">
              ID Card
            </SoftTypography>
          ),
          ciphertext: (
            <SoftTypography variant="caption" color="secondary" fontWeight="medium">
              {responseData.id_card_rc4}
            </SoftTypography>
          ),
        },
        {
          data: (
            <SoftTypography variant="caption" fontWeight="medium">
              Username
            </SoftTypography>
          ),
          ciphertext: (
            <SoftTypography variant="caption" color="secondary" fontWeight="medium">
              {responseData.username_rc4}
            </SoftTypography>
          ),
        },
        {
          data: (
            <SoftTypography variant="caption" fontWeight="medium">
              Password
            </SoftTypography>
          ),
          ciphertext: (
            <SoftTypography variant="caption" color="secondary" fontWeight="medium">
              {responseData.password_rc4}
            </SoftTypography>
          ),
        },
      ]);

      // Set DES rows
      setDesRows([
        {
          data: (
            <SoftTypography variant="caption" fontWeight="medium">
              Name
            </SoftTypography>
          ),
          ciphertext: (
            <SoftTypography variant="caption" color="secondary" fontWeight="medium">
              {responseData.name_des}
            </SoftTypography>
          ),
        },
        {
          data: (
            <SoftTypography variant="caption" fontWeight="medium">
              Email
            </SoftTypography>
          ),
          ciphertext: (
            <SoftTypography variant="caption" color="secondary" fontWeight="medium">
              {responseData.email_des}
            </SoftTypography>
          ),
        },
        {
          data: (
            <SoftTypography variant="caption" fontWeight="medium">
              Phone
            </SoftTypography>
          ),
          ciphertext: (
            <SoftTypography variant="caption" color="secondary" fontWeight="medium">
              {responseData.notelp_des}
            </SoftTypography>
          ),
        },
        {
          data: (
            <SoftTypography variant="caption" fontWeight="medium">
              Address
            </SoftTypography>
          ),
          ciphertext: (
            <SoftTypography variant="caption" color="secondary" fontWeight="medium">
              {responseData.address_des}
            </SoftTypography>
          ),
        },
        {
          data: (
            <SoftTypography variant="caption" fontWeight="medium">
              ID Card
            </SoftTypography>
          ),
          ciphertext: (
            <SoftTypography variant="caption" color="secondary" fontWeight="medium">
              {responseData.id_card_des}
            </SoftTypography>
          ),
        },
        {
          data: (
            <SoftTypography variant="caption" fontWeight="medium">
              Username
            </SoftTypography>
          ),
          ciphertext: (
            <SoftTypography variant="caption" color="secondary" fontWeight="medium">
              {responseData.username_des}
            </SoftTypography>
          ),
        },
        {
          data: (
            <SoftTypography variant="caption" fontWeight="medium">
              Password
            </SoftTypography>
          ),
          ciphertext: (
            <SoftTypography variant="caption" color="secondary" fontWeight="medium">
              {responseData.password_des}
            </SoftTypography>
          ),
        },
      ]);
    } catch (error) {
      console.error("Error fetching data:", error);
    }
  };

  useEffect(() => {
    fetchData(); // Call the fetch function on component load
  }, []);

  const columns = [
    { name: "data", align: "center", width: "20%" },
    { name: "ciphertext", align: "left", width: "80%" },
  ];

  return (
    <DashboardLayout>
      <DashboardNavbar />
      <SoftBox py={3}>
        <SoftBox mb={3}>
          <Card>
            <SoftBox display="flex" justifyContent="space-between" alignItems="center" p={3}>
              <SoftTypography variant="h6">AES Encryption</SoftTypography>
            </SoftBox>
            <SoftBox
              sx={{
                "& .MuiTableRow-root:not(:last-child)": {
                  "& td": {
                    borderBottom: ({ borders: { borderWidth, borderColor } }) =>
                      `${borderWidth[1]} solid ${borderColor}`,
                  },
                },
              }}
            >
              <Table columns={columns} rows={aesRows} />
            </SoftBox>
          </Card>
        </SoftBox>
        <SoftBox mb={3}>
          <Card>
            <SoftBox display="flex" justifyContent="space-between" alignItems="center" p={3}>
              <SoftTypography variant="h6">RC4 Encryption</SoftTypography>
            </SoftBox>
            <SoftBox
              sx={{
                "& .MuiTableRow-root:not(:last-child)": {
                  "& td": {
                    borderBottom: ({ borders: { borderWidth, borderColor } }) =>
                      `${borderWidth[1]} solid ${borderColor}`,
                  },
                },
              }}
            >
              <Table columns={columns} rows={rc4Rows} />
            </SoftBox>
          </Card>
        </SoftBox>
        <Card>
          <SoftBox display="flex" justifyContent="space-between" alignItems="center" p={3}>
            <SoftTypography variant="h6">DES Encryption</SoftTypography>
          </SoftBox>
          <SoftBox
            sx={{
              "& .MuiTableRow-root:not(:last-child)": {
                "& td": {
                  borderBottom: ({ borders: { borderWidth, borderColor } }) =>
                    `${borderWidth[1]} solid ${borderColor}`,
                },
              },
            }}
          >
            <Table columns={columns} rows={desRows} />
          </SoftBox>
        </Card>
      </SoftBox>
      <Footer />
    </DashboardLayout>
  );
}

export default Encrypted;
