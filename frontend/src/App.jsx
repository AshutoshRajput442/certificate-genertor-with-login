// import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
// import LoginSignup from "./LoginSignup";
// import CourseList from "./CourseList";
// import SelectedCourse from "./SelectedCourse";
// import CertificateGenerator from "./CertificateGenerator"; // ✅ New certificate component

// const App = () => {
//   return (
//     <Router>
//       <Routes>
//         <Route path="/" element={<LoginSignup />} />
//         <Route path="/courses" element={<CourseList />} />
//         <Route path="/select-course" element={<CourseList />} />
//         <Route path="/selected-courses/:studentId" element={<SelectedCourse />} />
     
//         <Route path="/certificate/:studentId/:courseId" element={<CertificateGenerator />} />

     
     
//       </Routes>
//     </Router>
//   );
// };

// export default App;
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import LoginSignup from "./LoginSignup";
import CourseList from "./CourseList";
import SelectedCourse from "./SelectedCourse";
import CertificateGenerator from "./CertificateGenerator"; // ✅ Certificate Component

const App = () => {
  return (
    <Router>
      <Routes>
        {/* Authentication Route */}
        <Route path="/" element={<LoginSignup />} />

        {/* Course Selection Routes */}
        <Route path="/courses" element={<CourseList />} />
        <Route path="/select-course" element={<CourseList />} /> {/* ✅ Kept for consistency */}

        {/* Selected Courses for a Student */}
        <Route path="/selected-courses/:studentId" element={<SelectedCourse />} />

        {/* Certificate Generation */}
        <Route path="/generate-certificate/:studentId/:courseId" element={<CertificateGenerator />} />
        </Routes>
    </Router>
  );
};

export default App;
