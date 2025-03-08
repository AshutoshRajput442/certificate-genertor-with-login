import { useParams } from "react-router-dom";
import { useEffect, useState } from "react";

const CertificateGenerator = () => {
  const { courseId } = useParams();
  const [message, setMessage] = useState("");

  const generateCertificate = async () => {
    const token = localStorage.getItem("token");

    try {
      const response = await fetch(`http://localhost:8080/generate-certificate/${courseId}`, {
        method: "POST",
        headers: { Authorization: token },
      });

      const data = await response.json();
      if (response.ok) {
        setMessage("Certificate generated successfully!");
      } else {
        setMessage(data.error || "Failed to generate certificate");
      }
    } catch {
      setMessage("Server error");
    }
  };

  useEffect(() => {
    generateCertificate();
  }, []);

  return (
    <div className="certificate-container">
      <h2>Certificate Status</h2>
      {message && <p className="message">{message}</p>}
    </div>
  );
};

export default CertificateGenerator;
