import React, { useState, useEffect, useMemo } from 'react';

function ErrorPage() {
  const [colorIndex, setColorIndex] = useState(0);
  const errorMessage = '404 Not Found';

  const colors = useMemo(
    () => ["#FF0000", "#FF7F00", "#FFFF00", "#00FF00", "#0000FF", "#4B0082", "#9400D3"],
    []
  );

  const textColors = useMemo(() => {
    return errorMessage.split('').map((char, index) => ({
      char,
      color: colors[(index + colorIndex) % colors.length],
    }));
  }, [colorIndex, colors]);

  useEffect(() => {
    const interval = setInterval(() => {
      setColorIndex((prevIndex) => (prevIndex + 1) % colors.length);
    }, 150);

    return () => clearInterval(interval);
  }, [colors]);
  return (
    <div className="error-page">
      <h1 className="error-heading">
        {textColors.map(({ char }, index) => (
          <span key={index} style={{ color: colors[(index + colorIndex) % colors.length] }}>
            {char}
          </span>
        ))}
      </h1>
      <p className="back-link">
        Back to <a href="/">Home</a>
      </p>
    </div>
  );
};

export default ErrorPage;

