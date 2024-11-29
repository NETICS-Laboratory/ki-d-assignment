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
import { useNavigate } from "react-router-dom";

function Inbox() {
  const [requests, setRequests] = useState([]);

  // Fetch data function
  const fetchData = async () => {
    try {
      const token = localStorage.getItem("token");

      const response = await axios.get("http://127.0.0.1:8090/api/user/request-access", {
        headers: {
          Authorization: `Bearer ${token}`, // Add the Bearer token here
        },
        params: {
          type: "received", // Add query parameters here
        },
      });

      setRequests(response.data.data);
    } catch (error) {
      console.error("Error fetching data:", error);
    }
  };

  useEffect(() => {
    fetchData(); // Call the fetch function on component load
  }, []);

  const handleApprove = (id) => {
    const token = localStorage.getItem("token");

    axios
      .put(
        `http://127.0.0.1:8090/api/user/request-access/${id}`,
        { status: "approved" }, // Request body
        {
          headers: {
            Authorization: `Bearer ${token}`, // Authorization header
          },
        }
      )
      .then((response) => {
        if (response.data.status) {
          // Handle success (e.g., update the list, show a success message)
          window.location.reload();
        }
      })
      .catch((error) => {
        console.error("Error approving request:", error);
      });
  };

  const handleDeny = (id) => {
    const token = localStorage.getItem("token");

    axios
      .put(
        `http://127.0.0.1:8090/api/user/request-access/${id}`,
        { status: "denied" }, // Request body
        {
          headers: {
            Authorization: `Bearer ${token}`, // Authorization header
          },
        }
      )
      .then((response) => {
        if (response.data.status) {
          // Handle success (e.g., update the list, show a success message)
          window.location.reload();
        }
      })
      .catch((error) => {
        console.error("Error denying request:", error);
      });
  };

  return (
    <DashboardLayout>
      <DashboardNavbar />
      <SoftBox py={3}>
        {requests.map((request) => (
          <SoftBox mb={2} key={request.id}>
            <Card>
              <SoftBox
                sx={{
                  display: "flex",
                  justifyContent: "space-between",
                  p: 3,
                }}
              >
                <SoftBox alignContent="center">
                  <SoftButton variant="gradient" color="info">
                    User {request.user_id}
                  </SoftButton>
                  <SoftButton
                    variant="gradient"
                    color="light"
                    sx={{
                      ml: 1,
                    }}
                  >
                    {request.status}
                  </SoftButton>
                </SoftBox>

                <SoftBox alignContent="center">
                  <SoftButton
                    variant="gradient"
                    color="success"
                    sx={{
                      mx: 1,
                    }}
                    onClick={() => handleApprove(request.id)}
                  >
                    Approve
                  </SoftButton>
                  <SoftButton
                    variant="gradient"
                    color="warning"
                    onClick={() => handleDeny(request.id)}
                  >
                    Deny
                  </SoftButton>
                </SoftBox>
              </SoftBox>
            </Card>
          </SoftBox>
        ))}
      </SoftBox>
    </DashboardLayout>
  );
}

export default Inbox;
