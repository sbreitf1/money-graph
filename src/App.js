import { Container } from '@mui/material';
import { BrowserRouter, Route, Routes } from 'react-router-dom';

import Dashboard from './Dashboard';

export default function App() {
  return (
    <Container component="main" maxWidth="sm">
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<Dashboard />} />
        </Routes>
      </BrowserRouter>
    </Container>
  );
}
