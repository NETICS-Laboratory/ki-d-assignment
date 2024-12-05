# KI-D Installing OpenXPKI Assignment

## **1️⃣ What is OpenXPKI?**

OpenXPKI is an **open-source Public Key Infrastructure (PKI) platform** designed to manage and automate the  **digital certificates**. These certificates are essential for establishing trust and enabling secure communication in digital environments, such as websites, applications, and devices. OpenXPKI provides the tools to create, validate, and manage certificates, ensuring robust security across a wide range of use cases.


## **2️⃣ Roles and Their Responsibilities**

OpenXPKI encompasses three primary roles: **Certification Authority (CA)**, **Registration Authority (RA)**, and **Common User (CU)**. Each role has distinct responsibilities that contribute to the integrity and functionality of the PKI system.

### 1. Certification Authority (CA)
The **CA** is the cornerstone of the PKI system, responsible for issuing and managing digital certificates. It serves as the ultimate source of trust.

- **Responsibilities:**
  - **Certificate Issuance:** Creates digital certificates after receiving verified requests.
  - **Certificate Signing:** Signs certificates with its private key to establish their authenticity.
  - **Certificate Revocation:** Invalidates certificates that are compromised or no longer required.
  - **Maintaining CRL (Certificate Revocation List):** Publishes a list of revoked certificates to prevent misuse.

### 2. Registration Authority (RA)
The **RA** serves as the intermediary between the CA and certificate requesters, validating identities before certificates are issued.

- **Responsibilities:**
  - **Identity Verification:** Confirms the legitimacy of entities requesting certificates.
  - **Request Approval:** Ensures that only authenticated and authorized requests are forwarded to the CA.
  - **Request Forwarding:** Relays approved requests to the CA for certificate generation.

### 3. Common User
The **Common User** represents the end entity that requests and utilizes the certificates for secure operations, such as website encryption or device authentication.

- **Responsibilities:**
  - **Request Certificates:** Initiates requests for digital certificates via the RA.
  - **Implement Certificates:** Deploys the issued certificates for intended purposes (e.g., HTTPS encryption, email security).
  - **Monitor Certificates:** Tracks certificate validity and requests renewals as needed.


## **3️⃣ The Workflow of OpenXPKI88**

The roles in OpenXPKI collaborate systematically to deliver a secure and efficient PKI solution. Here is a typical workflow:

1. **Certificate Request:**
   - The Common User (a website owner) submits a request for a certificate to the RA.

2. **Identity Validation:**
   - The RA authenticates the request by verifying the identity and legitimacy of the requester.

3. **Certificate Issuance:**
   - Upon approval, the RA forwards the request to the CA, which generates and digitally signs the certificate.

4. **Certificate Delivery:**
   - The RA securely provides the issued certificate to the Common User.

5. **Deployment and Use:**
   - The Common User implements the certificate in their environment (a web server for HTTPS).


## **5️⃣ Installing Docker Desktop on macOS from the Installer Package**
#### **Step 1: Download Docker Desktop**
1. Open your web browser and visit the official Docker website at https://www.docker.com/products/docker-desktop/.
2. On the Docker Desktop for Mac page, click on the “Download for Mac” button. If you have an Apple Silicon Mac (M1 or M2 chip), make sure to select the “Apple Chip” option.
3. The Docker Desktop installer package (.dmg file) will start downloading. Once the download is complete, proceed to the next step.

#### **Step 2: Install Docker Desktop**
1. Locate the downloaded Docker Desktop installer package (.dmg file) in your Downloads folder or the location you specified during the download.
2. In the mounted disk image, you will see the Docker icon and an Applications folder shortcut. Click and drag the Docker icon onto the Applications folder shortcut. This will copy the Docker Desktop application to your Mac’s Applications folder.

![drag app](https://github.com/laurivasyyy/PBKK-Tugas-2-WebCaptureApp/blob/e0e13ca8548c5423ad342cce1cef3fcc7cbc996b/img/Screenshot%202024-12-04%20at%2013.08.55.png)

#### **Step 3: Verify Docker Desktop Installation**
1. Once Docker Desktop finishes initializing, you will see the Docker Desktop user interface.
2. To further verify the installation, open the Terminal and run the following command:

```bash
docker --version
```
```bash
docker-compose --version
```
![version](https://github.com/laurivasyyy/PBKK-Tugas-2-WebCaptureApp/blob/e0e13ca8548c5423ad342cce1cef3fcc7cbc996b/img/Screenshot%202024-12-04%20at%2013.17.04.png)

## **6️⃣ OpenXPKI Installation**
We install OpenXPKI by first cloning the repository of OpenXPKI. We can do that by running this code:

```bash
git clone https://github.com/openxpki/openxpki-docker.git
```

Now that the cloning process is done, make sure to go to the directory of the repo that we just clone.

Next while we are inside the OpenXPKI directory, we clone the repository for the config. The command to do that is as follows:

```bash
git clone https://github.com/openxpki/openxpki-config.git \
 --single-branch --branch=community
 ```

CAUTIONS: To prevent the server from crashing when the database isn’t working, it’s suggested to duplicate the configuration into the local.yaml file. You can do that by using the command as follows:

```
cp contrib/wait_on_init.yaml  openxpki-config/config.d/system/local.yaml
```

 After all that is done, inside the directory should look like this:

![image](https://github.com/laurivasyyy/PBKK-Tugas-2-WebCaptureApp/blob/e0e13ca8548c5423ad342cce1cef3fcc7cbc996b/img/Screenshot%202024-12-04%20at%2013.30.31.png)

Now, to run the docker-compose. Use below make command to start

```bash
brew install make
make compose
```

![image](https://github.com/laurivasyyy/PBKK-Tugas-2-WebCaptureApp/blob/e0e13ca8548c5423ad342cce1cef3fcc7cbc996b/img/Screenshot%202024-12-04%20at%2013.19.48.png)

or

you can just start your docker compose directly from the docker desktop

![image](https://github.com/laurivasyyy/PBKK-Tugas-2-WebCaptureApp/blob/e0e13ca8548c5423ad342cce1cef3fcc7cbc996b/img/Screenshot%202024-12-04%20at%2013.36.32.png)

The Web-Server is now started, to Open the OpenXPKI Web, you can access https://localhost:8443/


## **6️⃣ Using OpenXPKI as Common User (CU)**

#### **Step 1: Login using the Test Account**
![image](https://github.com/laurivasyyy/PBKK-Tugas-2-WebCaptureApp/blob/1e9d4e2d4765300ffb7baf414f7cbd46795f4f81/img/Screenshot%202024-12-04%20at%2013.51.48.png)

You can choose to use the `Test Account` and enter `alice` or `bob` as the username for a common user, the password for all accounts is `openxpki`.

#### **Step 2: Create a Certificate Request**

To create a profile, we selected `TLS/Web Server` as the certificate profile and the subject style was left as the default.

For the request type, we chose `Generate Key on PKI` to allow the server to generate the private key for you, which is a much safer option.

![image](https://github.com/laurivasyyy/PBKK-Tugas-2-WebCaptureApp/blob/1e9d4e2d4765300ffb7baf414f7cbd46795f4f81/img/Screenshot%202024-12-04%20at%2020.38.14.png)

To generate the key on PKI, we must select the key parameters. Here, we are choosing `RSA` with a key length of `2048 bit`.
![image](https://github.com/laurivasyyy/PBKK-Tugas-2-WebCaptureApp/blob/1e9d4e2d4765300ffb7baf414f7cbd46795f4f81/img/Screenshot%202024-12-04%20at%2020.38.18.png)

Now, we need to edit the item in the main subject of the certification request you can enter the `Hostname`, `Additional Hostname`, and `Application Name`.
![image](https://github.com/laurivasyyy/PBKK-Tugas-2-WebCaptureApp/blob/1e9d4e2d4765300ffb7baf414f7cbd46795f4f81/img/Screenshot%202024-12-04%20at%2020.38.32.png)

To assist our registration officers when validating the certificate request we need to edit the certificate info, in here all the input fields was left as the default.
![image](https://github.com/laurivasyyy/PBKK-Tugas-2-WebCaptureApp/blob/1e9d4e2d4765300ffb7baf414f7cbd46795f4f81/img/Screenshot%202024-12-04%20at%2020.38.37.png)

Once you’ve completed all the data inputs, you can review your certificate and edit the subject if needed. If everything looks correct, click Continue to submit the request.

![image](https://github.com/laurivasyyy/PBKK-Tugas-2-WebCaptureApp/blob/1e9d4e2d4765300ffb7baf414f7cbd46795f4f81/img/Screenshot%202024-12-04%20at%2020.38.54.png)


After submitting your request, the server will automatically generate a password for you. This password will be used later to download and install your private key.

![image](https://github.com/laurivasyyy/PBKK-Tugas-2-WebCaptureApp/blob/1e9d4e2d4765300ffb7baf414f7cbd46795f4f81/img/Screenshot%202024-12-04%20at%2020.39.19.png)

Once your request is successfully submitted, you will need to wait for approval by the RA (Registration Authority).

![image](https://github.com/laurivasyyy/PBKK-Tugas-2-WebCaptureApp/blob/1e9d4e2d4765300ffb7baf414f7cbd46795f4f81/img/Screenshot%202024-12-04%20at%2021.08.54.png)

Once the certification request already issued by the registration authority you can do this :

1. **Downloading the private key :** 
You can go to the `Action` section to `download the private key` by entering the password you received when you created the request.

Once you succeed, you will be directed to the keystore download button.
![image](https://github.com/laurivasyyy/PBKK-Tugas-2-WebCaptureApp/blob/1e9d4e2d4765300ffb7baf414f7cbd46795f4f81/img/Screenshot%202024-12-04%20at%2021.12.13.png)

2. **Revoking the Certificate :** 
You can go to the navbar section `Revoke Certificate`, then enter the `Certificate Identifier` and choose `Certificate Authority Key was compromised` as the reason code.
![image](https://github.com/laurivasyyy/PBKK-Tugas-2-WebCaptureApp/blob/1e9d4e2d4765300ffb7baf414f7cbd46795f4f81/img/Screenshot%202024-12-04%20at%2021.12.35.png)

You can see the summary of your revocation request below. After that you can press the submit button to send the revocation request of this certificate.

![image](https://github.com/laurivasyyy/PBKK-Tugas-2-WebCaptureApp/blob/1e9d4e2d4765300ffb7baf414f7cbd46795f4f81/img/Screenshot%202024-12-04%20at%2021.13.08.png
)

## **7️⃣ Using OpenXPKI as Registration Authority (RA)**

#### **Step 1: Login using the Test Account**
![image](https://github.com/laurivasyyy/PBKK-Tugas-2-WebCaptureApp/blob/1e9d4e2d4765300ffb7baf414f7cbd46795f4f81/img/Screenshot%202024-12-04%20at%2013.51.48.png)

You can choose to use the `Test Account` and enter `rob` as the username for the registration authority, the password for all accounts is `openxpki`.

#### **Step 2: Managing the Certificate Request**
To manage Certificate Request, you can click Home  
-> My Task. Here you can approve or revoke any 
certificate requests from the common users.

![image](https://github.com/laurivasyyy/PBKK-Tugas-2-WebCaptureApp/blob/1e9d4e2d4765300ffb7baf414f7cbd46795f4f81/img/Screenshot%202024-12-04%20at%2021.10.03.png)

You will see the common user's certificate request pending approval. You can click on any of them to proceed with approving the request.

![image](https://github.com/laurivasyyy/PBKK-Tugas-2-WebCaptureApp/blob/1e9d4e2d4765300ffb7baf414f7cbd46795f4f81/img/Screenshot%202024-12-04%20at%2021.10.29.png)

#### **Step 3: Validating the Certificate Request Status**

You must log out from the Registration Authority account and log back in using the Common User's account. Once logged in, navigate to Home -> My Certificate to check if the certificate request has been issued.

![image](https://github.com/laurivasyyy/PBKK-Tugas-2-WebCaptureApp/blob/1e9d4e2d4765300ffb7baf414f7cbd46795f4f81/img/Screenshot%202024-12-04%20at%2021.11.28.png)

## **8️⃣ Using OpenXPKI as Certification Authority (CA)**

#### **Step 1: Login using the Test Account**
![image](https://github.com/laurivasyyy/PBKK-Tugas-2-WebCaptureApp/blob/1e9d4e2d4765300ffb7baf414f7cbd46795f4f81/img/Screenshot%202024-12-04%20at%2013.51.48.png)

You can choose to use the `Test Account` and enter `rob` as the username for the certification authority, the password for all accounts is `openxpki`. By default, Rob has a dual role that also acts as a CA.

#### **Step 2: Trigger CRL Issue**
A revocation list is created when there are new revocations or if the current list is near its expiry. But you can always force to do so by ticking all of the checkbox

![image](https://github.com/laurivasyyy/PBKK-Tugas-2-WebCaptureApp/blob/1e9d4e2d4765300ffb7baf414f7cbd46795f4f81/img/Screenshot%202024-12-04%20at%2021.14.50.png)

![image](https://github.com/laurivasyyy/PBKK-Tugas-2-WebCaptureApp/blob/1e9d4e2d4765300ffb7baf414f7cbd46795f4f81/img/Screenshot%202024-12-04%20at%2021.16.03.png)

here you can click force wake up and the revocation lists item will be added.










