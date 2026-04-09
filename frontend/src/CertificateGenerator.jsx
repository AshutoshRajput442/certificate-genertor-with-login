import { useState } from "react";
import { useParams } from "react-router-dom";

const GenerateCertificate = () => {
  const { studentId, courseId } = useParams();
  const [loading, setLoading] = useState(false);
  const [certificateData, setCertificateData] = useState(null);
  const [error, setError] = useState(null);

  const handleGenerateCertificate = async () => {
    setLoading(true);
    setError(null);

    try {
      const response = await fetch(`http://localhost:8080/generate-certificate/${studentId}/${courseId}`);
      const data = await response.json();

      if (response.ok) {
        setCertificateData(data);
        
        // Auto-download the PDF
        const downloadLink = document.createElement("a");
        downloadLink.href = `http://localhost:8080${data.pdf_url}`;
        downloadLink.download = `certificate_${studentId}_${courseId}.pdf`;
        document.body.appendChild(downloadLink);
        downloadLink.click();
        document.body.removeChild(downloadLink);
      } else {
        setError(data.error || "Failed to generate certificate");
      }
    } catch (err) {
      setError("Something went wrong. Please try again.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div>
      <h2>Generate Certificate</h2>
      <button onClick={handleGenerateCertificate} disabled={loading}>
        {loading ? "Generating..." : "Generate Certificate"}
      </button>
      
      {certificateData && (
        <div>
          <h3>Certificate Details</h3>
          <p><strong>Student ID:</strong> {certificateData.student_id}</p>
          <p><strong>Student Name:</strong> {certificateData.student_name}</p>
          <p><strong>Course Name:</strong> {certificateData.course_name}</p>
          <a href={`http://localhost:8080${certificateData.pdf_url}`} download>
            Download Certificate
          </a>
        </div>
      )}

      {error && <p style={{ color: "red" }}>{error}</p>}
    </div>
  );
};

export default GenerateCertificate;




