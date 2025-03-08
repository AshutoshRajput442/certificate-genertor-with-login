import { useState, useEffect } from "react";
import axios from "axios";
import "./selectedcourses.css";

const SelectedCourse = ({ studentId }) => {
  const [selectedCourses, setSelectedCourses] = useState([]);

  // API call to fetch selected courses
  const fetchSelectedCourses = async () => {
    try {
      const response = await axios.get(`http://localhost:8080/selected-courses/${studentId}`);
      setSelectedCourses(response.data);
    } catch (error) {
      console.error("Error fetching selected courses:", error);
    }
  };

  useEffect(() => {
    fetchSelectedCourses();
  }, [studentId]);

  return (
    <div>
      <h2>Your Selected Courses</h2>
      <button onClick={fetchSelectedCourses} className="fetch-btn">
        Show Your Courses
      </button>
      <ul>
        {selectedCourses.length > 0 ? (
          selectedCourses.map((course) => (
            <li key={course.id}>{course.name}</li>
          ))
        ) : (
          <p>No courses selected yet.</p>
        )}
      </ul>
    </div>
  );
};

export default SelectedCourse;



// // function SelectedCourse() {
// //   return (
// //     <h1>hello</h1>
// //   )
// // }

// // export default SelectedCourse;


// import { useState, useEffect } from "react";
// import { useParams } from "react-router-dom";

// const SelectedCourses = () => {
//   const { studentId } = useParams(); // ✅ Get student ID from URL
//   const [selectedCourses, setSelectedCourses] = useState([]);
//   const [loading, setLoading] = useState(true);

//   useEffect(() => {
//     fetch(`http://localhost:8080/selected-courses/${studentId}`)
//       .then((res) => res.json())
//       .then((data) => {
//         console.log("Fetched Selected Courses:", data); // ✅ Debugging
//         setSelectedCourses(data);
//         setLoading(false);
//       })
//       .catch((err) => {
//         console.error("Error fetching selected courses:", err);
//         setLoading(false);
//       });
//   }, [studentId]);

//   if (loading) return <p>Loading...</p>;

//   return (
//     <div>
//       <h2>Your Selected Courses</h2>
//       {selectedCourses.length === 0 ? (
//         <p>No courses selected yet.</p>
//       ) : (
//         <ul>
//           {selectedCourses.map((course) => (
//             <li key={course.id}>{course.name}</li>
//           ))}
//         </ul>
//       )}
//     </div>
//   );
// };

// export default SelectedCourses;
