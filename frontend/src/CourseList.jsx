import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import "./courselist.css";

const CourseList = () => {
  const [courses, setCourses] = useState([]);
  const [message, setMessage] = useState("");
  const navigate = useNavigate();

  // Retrieve logged-in student ID (Replace with actual auth logic)
  const studentId = localStorage.getItem("studentId") || 7; // Replace with actual authentication method

  // Fetching Courses from Backend
  useEffect(() => {
    fetch("http://localhost:8080/courses")
      .then((res) => res.json())
      .then((data) => setCourses(data))
      .catch((err) => console.error("Error fetching courses:", err));
  }, []);

  // Selecting Course and sending data to backend

  // particular course is selected by the  particular student  and send to selected_courses table

  // This function handleSelectCourse is responsible for sending a request to the backend when a student selects a course.

  // const handleSelectCourse = (courseId) => {
  //   fetch("http://localhost:8080/select-course", {
  //     // ✅ Fixed API URL
  //     method: "POST",
  //     headers: { "Content-Type": "application/json" },
  //     body: JSON.stringify({ course_id: courseId, student_id: studentId }), // ✅ Use dynamic student ID
  //   })
  //     .then((res) => res.json())
  //     .then((data) => {
  //       console.log("API Response:", data); // ✅ Debugging
  //       if (data.message) {
  //         setMessage(`✅ ${data.message}: ${courseName}`);
  //       } else {
  //         setMessage("❌ Failed to select course.");
  //       }
  //       setTimeout(() => setMessage(""), 3000);
  //     })
  //     .catch((err) => {
  //       console.error("Error selecting course:", err);
  //       setMessage("❌ Error selecting course.");
  //     });
  // };







  const handleSelectCourse = async (courseId) => {
    const studentId = localStorage.getItem("studentId");
  
    if (!studentId) {
      alert("Student ID not found! Please log in again.");
      return;
    }
  
    try {
      const response = await fetch("http://localhost:8080/select-course", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ student_id: parseInt(studentId), course_id: courseId }),
      });
  
      const data = await response.json();
      console.log("API Response:", data);
  
      if (!response.ok) {
        throw new Error(data.error || "Failed to select course.");
      }
  
      setMessage(`✅ ${data.message}`);
      setTimeout(() => setMessage(""), 3000);
    } catch (error) {
      console.error("Error selecting course:", error);
      setMessage("❌ " + error.message);
    }
  };
  

  return (
    <div className="course-list">
      <h2>Available Courses</h2>
      {message && <p className="success-message">{message}</p>}

      {/* ✅ Navigate to the selected courses page dynamically */}
      <button
        className="show-courses-btn"
        onClick={() => navigate(`/selected-courses/${studentId}`)}
      >
        Show Your Courses
      </button>

      <ul>
        {courses.length > 0 ? (
          courses.map((course) => (
            <li key={course.id}>
              {course.name}
              <button onClick={() => handleSelectCourse(course.id, course.name)}>Select</button>
            </li>
          ))
        ) : (
          <p>Loading courses...</p>
        )}
      </ul>
    </div>
  );
};

export default CourseList;

// function CourseList() {
//   return (
//     <h1>hello</h1>
//   )
// }

// export default CourseList;

// import { useState, useEffect } from "react";
// import { useNavigate } from "react-router-dom";
// import "./courselist.css";

// const CourseList = () => {
//   const [courses, setCourses] = useState([]);
//   const [message, setMessage] = useState("");
//   const navigate = useNavigate();

//   // Retrieve logged-in student ID (Replace with actual auth logic)
//   const studentId = localStorage.getItem("studentId") || 7; // Replace with actual authentication method

//   // Fetch Courses from Backend
//   useEffect(() => {
//     fetch("http://localhost:8080/courses")
//       .then((res) => res.json())
//       .then((data) => setCourses(data))
//       .catch((err) => console.error("Error fetching courses:", err));
//   }, []);

//   // Function to Select a Course
//   const handleSelectCourse = (courseId, courseName) => {
//     fetch("http://localhost:8080/select-course", { // ✅ Fixed API URL
//       method: "POST",
//       headers: { "Content-Type": "application/json" },
//       body: JSON.stringify({ course_id: courseId, student_id: studentId }),
//     })
//       .then((res) => res.json())
//       .then((data) => {
//         console.log("API Response:", data); // ✅ Debugging
//         if (data.message) {
//           setMessage(`✅ ${data.message}: ${courseName}`);
//         } else {
//           setMessage("❌ Failed to select course.");
//         }
//         setTimeout(() => setMessage(""), 3000);
//       })
//       .catch((err) => {
//         console.error("Error selecting course:", err);
//         setMessage("❌ Error selecting course.");
//       });
//   };

//   return (
//     <div className="course-list">
//       <h2>Available Courses</h2>
//       {message && <p className="success-message">{message}</p>}

//       {/* Navigate to the selected courses page dynamically */}
//       <button className="show-courses-btn" onClick={() => navigate(`/selected-courses/${studentId}`)}>
//         Show Your Courses
//       </button>

//       <ul>
//         {courses.length > 0 ? (
//           courses.map((course) => (
//             <li key={course.id}>
//               {course.name}
//               <button onClick={() => handleSelectCourse(course.id, course.name)}>Select</button>
//             </li>
//           ))
//         ) : (
//           <p>Loading courses...</p>
//         )}
//       </ul>
//     </div>
//   );
// };

// export default CourseList;
