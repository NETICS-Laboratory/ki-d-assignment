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

// Mock API Data
const apiResponse = {
  status: true,
  message: "Berhasil Mendapatkan User",
  errors: null,
  data: {
    id: "ca587cab-b1df-4fe1-91a1-f17d0e9ce91c",
    username: "test",
    name: "test",
    email: "test@test.com",
    no_telp: "12345",
    address: "test",
    id_card: "uploads/test/encrypted/cloud-139.png",
  },
};

function Decrypted() {
  const [data, setData] = useState({});
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [decryptedData, setDecryptedData] = useState(null);
  const [rows, setRows] = useState([]);

  const fetchData = async () => {
    try {
      const token = localStorage.getItem("token");

      const response = await axios.get("http://127.0.0.1:8090/api/user/me-decrypted", {
        headers: {
          Authorization: `Bearer ${token}`, // Add the Bearer token here
        },
      });

      const responseData = response.data.data;

      setData(responseData);

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

  const fetchDecryptedID = async () => {
    try {
      const token = localStorage.getItem("token");

      const response = await axios.get("http://127.0.0.1:8090/api/user/idcard-decrypted", {
        headers: {
          Authorization: `Bearer ${token}`, // Add the Bearer token here
        },
      });
    } catch (error) {
      console.error("Error fetching decrypted ID Card:", error);
    }
  };

  // Table columns setup
  const columns = [
    { name: "data", label: "Data", align: "center", width: "20%" },
    { name: "plaintext", label: "Plaintext", align: "left", width: "80%" },
  ];

  const handleLogin = async (e) => {
    e.preventDefault();

    try {
      // Call the login API
      const response = await axios.post("http://127.0.0.1:8090/api/user/login", {
        username,
        password,
      });

      // If the API response includes a token, the login is successful
      if (response.data.data.token) {
        setIsLoggedIn(true);

        // Fetch the decrypted data
        fetchData();
        fetchDecryptedID();
      } else {
        alert("Invalid username or password");
      }
    } catch (error) {
      console.error("Error during login:", error);
      alert("Login failed. Please try again.");
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
                  Decrypt Data
                </SoftTypography>
                <SoftTypography variant="button" fontWeight="regular">
                  Fill in your credentials to decrypt your data
                </SoftTypography>
              </SoftBox>
              <Separator />
              <SoftBox pt={2} pb={3} px={3}>
                <SoftBox component="form" role="form" onSubmit={handleLogin}>
                  <SoftBox mb={2}>
                    <SoftInput
                      type="text"
                      name="username"
                      placeholder="Username"
                      value={username}
                      onChange={(e) => setUsername(e.target.value)}
                    />
                  </SoftBox>
                  <SoftBox mb={2}>
                    <SoftInput
                      type="password"
                      name="password"
                      placeholder="Password"
                      value={password}
                      onChange={(e) => setPassword(e.target.value)}
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

export default Decrypted;
