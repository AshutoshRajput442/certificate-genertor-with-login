import { useState, useEffect } from "react";
import { useParams } from "react-router-dom";

const SelectedCourses = () => {
  const { studentId } = useParams();
  const [selectedCourses, setSelectedCourses] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetch(`http://localhost:8080/selected-courses/${studentId}`)
      .then((res) => res.json())
      .then((data) => {
        console.log("Fetched Selected Courses:", data);
        setSelectedCourses(data);
        setLoading(false);
      })
      .catch((err) => {
        console.error("Error fetching selected courses:", err);
        setLoading(false);
      });
  }, [studentId]);

  const generateCertificate = async (courseId) => {
    try {
      const response = await fetch(
        `http://localhost:8080/generate-certificate/${studentId}/${courseId}`
      );
      const data = await response.json();

      if (response.ok) {
        console.log("Certificate Generated:", data);
        if (data.pdf_url) {
          // ✅ Auto-download the generated PDF
          const link = document.createElement("a");
          link.href = data.pdf_url;
          link.download = `certificate_${studentId}_${courseId}.pdf`;
          document.body.appendChild(link);
          link.click();
          document.body.removeChild(link);
        } else {
          alert("Certificate generation failed!");
        }
      } else {
        alert(data.error || "Failed to generate certificate.");
      }
    } catch (error) {
      console.error("Error generating certificate:", error);
      alert("Error generating certificate.");
    }
  };

  if (loading) return <p>Loading...</p>;

  return (
    <div>
      <h2>Your Selected Courses</h2>
      {selectedCourses.length === 0 ? (
        <p>No courses selected yet.</p>
      ) : (
        <ul>
          {selectedCourses.map((course) => (
            <li key={course.id}>
              {course.name}
              <button onClick={() => generateCertificate(course.id)}>
                Generate Certificate
              </button>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
};

export default SelectedCourses;
