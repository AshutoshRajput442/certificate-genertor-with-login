// import React, { useState, useEffect } from "react";
// import { useParams } from "react-router-dom";

// const CertificateGenerator = () => {
//   const { studentId, courseId } = useParams();
//   const [studentName, setStudentName] = useState("");
//   const [courseName, setCourseName] = useState("");
//   const [pdfUrl, setPdfUrl] = useState(null);
//   const [loading, setLoading] = useState(true);

//   // Fetch student details and course name
//   useEffect(() => {
//     const fetchStudentAndCourse = async () => {
//       try {
//         const response = await fetch(
//           `http://localhost:8080/certificate/${studentId}/${courseId}`
//         );
//         const result = await response.json();

//         if (response.ok) {
//           setStudentName(result.student);
//           setCourseName(result.course);
//           if (result.pdf_url) {
//             setPdfUrl(result.pdf_url);
//           }
//         } else {
//           alert(result.error || "Failed to fetch details.");
//         }
//       } catch (error) {
//         console.error("Error fetching data:", error);
//         alert("Server error: " + error.message);
//       } finally {
//         setLoading(false);
//       }
//     };

//     fetchStudentAndCourse();
//   }, [studentId, courseId]);

//   // Function to generate certificate
//   const handleGenerateCertificate = async () => {
//     try {
//       const response = await fetch(
//         `http://localhost:8080/certificate/${studentId}/${courseId}`,
//         {
//           method: "GET",
//           headers: { "Content-Type": "application/json" },
//         }
//       );

//       const result = await response.json();
//       if (response.ok) {
//         setPdfUrl(result.pdf_url);
//       } else {
//         alert(result.error || "Something went wrong!");
//       }
//     } catch (error) {
//       console.error("Error generating certificate:", error);
//       alert("Server error: " + error.message);
//     }
//   };

//   if (loading) return <p>Loading...</p>;

//   return (
//     <div>
//       <h1>Generate Certificate</h1>
//       <p>
//         <strong>Student ID:</strong> {studentId}
//       </p>
//       <p>
//         <strong>Student Name:</strong> {studentName}
//       </p>
//       <p>
//         <strong>Course ID:</strong> {courseId}
//       </p>
//       <p>
//         <strong>Course Name:</strong> {courseName}
//       </p>
//       {/* <p>Click the button below to generate your certificate:</p> */}

//       {/* <button>
//         <a href={`http://localhost:8080${pdfUrl}`} download>
//           Download Your Certificate
//         </a>
//       </button> */}

//       {/* Show PDF download link if available */}
//     </div>
//   );
// };

// export default CertificateGenerator;

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
