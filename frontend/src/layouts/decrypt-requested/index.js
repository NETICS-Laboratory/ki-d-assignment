import React, { useState, useEffect } from "react";
import Card from "@mui/material/Card";
import Button from "@mui/material/Button";
import TextField from "@mui/material/TextField";

// Soft UI Dashboard React components
import SoftBox from "components/SoftBox";
import SoftTypography from "components/SoftTypography";
import SoftInput from "components/SoftInput";
import SoftButton from "components/SoftButton";

// Authentication layout components
import Separator from "layouts/authentication/components/Separator";

// Soft UI Dashboard React examples
import DashboardLayout from "examples/LayoutContainers/DashboardLayout";
import DashboardNavbar from "examples/Navbars/DashboardNavbar";
import Footer from "examples/Footer";
import Table from "examples/Tables/Table";

import axios from "axios";

function DecryptedRequested() {
  const [data, setData] = useState({});
  // const [username, setUsername] = useState("");
  // const [password, setPassword] = useState("");

  const [requested_user_username, setRequested_user_username] = useState("");
  const [encrypted_key, setEncrypted_key] = useState("");
  const [encrypted_key_8_byte, setEncrypted_key_8_byte] = useState("");
  const [secret_key, setSecret_key] = useState("");
  const [secret_key_8_byte, setSecret_key_8_byte] = useState("");

  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [decryptedData, setDecryptedData] = useState(null);
  const [rows, setRows] = useState([]);

  const fetchRequestedUserData = async () => {
    try {
      const token = localStorage.getItem("token");

      const response = await axios.post(
        "http://127.0.0.1:8090/api/user/get-requested-user-data",
        {
          requested_user_username,
          secret_key,
          secret_key_8_byte,
        },
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );

      const responseData = response.data.data.requested_user;

      setData(responseData);
      setIsLoggedIn(true);

      // Set rows
      setRows([
        {
          data: (
            <SoftTypography variant="caption" fontWeight="medium">
              Name
            </SoftTypography>
          ),
          plaintext: (
            <SoftTypography variant="caption" color="secondary" fontWeight="medium">
              {responseData.name}
            </SoftTypography>
          ),
        },
        {
          data: (
            <SoftTypography variant="caption" fontWeight="medium">
              Email
            </SoftTypography>
          ),
          plaintext: (
            <SoftTypography variant="caption" color="secondary" fontWeight="medium">
              {responseData.email}
            </SoftTypography>
          ),
        },
        {
          data: (
            <SoftTypography variant="caption" fontWeight="medium">
              Phone
            </SoftTypography>
          ),
          plaintext: (
            <SoftTypography variant="caption" color="secondary" fontWeight="medium">
              {responseData.no_telp}
            </SoftTypography>
          ),
        },
        {
          data: (
            <SoftTypography variant="caption" fontWeight="medium">
              Address
            </SoftTypography>
          ),
          plaintext: (
            <SoftTypography variant="caption" color="secondary" fontWeight="medium">
              {responseData.address}
            </SoftTypography>
          ),
        },
        {
          data: (
            <SoftTypography variant="caption" fontWeight="medium">
              ID Card
            </SoftTypography>
          ),
          plaintext: (
            // <SoftTypography variant="caption" color="secondary" fontWeight="medium">
            //   {responseData.id_card}
            // </SoftTypography>

            <div>
              <img
                src={`${process.env.PUBLIC_URL}/idcard-example.jpg`} // Adjust the URL as necessary
                alt="ID Card"
                style={{ width: "100px", height: "auto" }} // Set desired size
              />
            </div>
          ),
        },
      ]);
    } catch (error) {
      console.error("Error fetching data:", error);
    }
  };

  // Table columns setup
  const columns = [
    { name: "data", label: "Data", align: "center", width: "20%" },
    { name: "plaintext", label: "Plaintext", align: "left", width: "80%" },
  ];

  const handleDecryptRequested = async (e) => {
    e.preventDefault();

    const token = localStorage.getItem("token");

    try {
      // Call the login API
      const response = await axios.post(
        "http://127.0.0.1:8090/api/user/decrypt-key",
        {
          encrypted_key,
          encrypted_key_8_byte,
        },
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );

      const responseData = response.data.data;

      setSecret_key(responseData.decrypted_key);
      setSecret_key_8_byte(responseData.decrypted_key_8_byte);

      fetchRequestedUserData();
    } catch (error) {
      console.error("Error during decrypting requested data:", error);
    }
  };

  return (
    <DashboardLayout>
      <DashboardNavbar />
      <SoftBox py={3}>
        <SoftBox mb={3}>
          {!isLoggedIn ? (
            <Card>
              <SoftBox p={3} textAlign="center">
                <SoftTypography variant="h5" fontWeight="medium">
                  Decrypt Requested Data
                </SoftTypography>
                <SoftTypography variant="button" fontWeight="regular">
                  Fill in the credentials to decrypt requested user&apos;s data
                </SoftTypography>
              </SoftBox>
              <Separator />
              <SoftBox pt={2} pb={3} px={3}>
                <SoftBox component="form" role="form" onSubmit={handleDecryptRequested}>
                  <SoftBox mb={2}>
                    <SoftInput
                      type="text"
                      name="requested_user_username"
                      placeholder="Requested Username"
                      value={requested_user_username}
                      onChange={(e) => setRequested_user_username(e.target.value)}
                    />
                  </SoftBox>
                  <SoftBox mb={2}>
                    <SoftInput
                      type="text"
                      name="encrypted_key"
                      placeholder="Encrypted Key"
                      value={encrypted_key}
                      onChange={(e) => setEncrypted_key(e.target.value)}
                    />
                  </SoftBox>
                  <SoftBox mb={2}>
                    <SoftInput
                      type="text"
                      name="encrypted_key_8_byte"
                      placeholder="Encrypted Key 8 Byte"
                      value={encrypted_key_8_byte}
                      onChange={(e) => setEncrypted_key_8_byte(e.target.value)}
                    />
                  </SoftBox>
                  <SoftBox mt={4} mb={1}>
                    <SoftButton type="submit" variant="gradient" color="dark" fullWidth>
                      Decrypt
                    </SoftButton>
                  </SoftBox>
                </SoftBox>
              </SoftBox>
            </Card>
          ) : (
            <SoftBox>
              <Card>
                <SoftBox display="flex" justifyContent="space-between" alignItems="center" p={3}>
                  <SoftTypography variant="h6">Decrypted Data</SoftTypography>
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
                  <Table columns={columns} rows={rows} />
                </SoftBox>
              </Card>
            </SoftBox>
          )}
        </SoftBox>
      </SoftBox>
    </DashboardLayout>
  );
}

export default DecryptedRequested;
