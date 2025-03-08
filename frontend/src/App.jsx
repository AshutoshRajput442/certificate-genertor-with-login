import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import LoginSignup from "./LoginSignup";
import CourseList from "./CourseList";
import SelectedCourse from "./SelectedCourse";
// import CertificateGenerator from "./CertificateGenerator";

const App = () => {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<LoginSignup />} />
        <Route path="/courses" element={<CourseList />} />
        <Route path="/select-course" element={<CourseList />} />
        <Route path="/selected-courses/:studentId" element={<SelectedCourse />} />
        {/* <Route path="/certificate/:courseId" element={<CertificateGenerator />} /> */}
      </Routes>
    </Router>
  );
};

export default App;
