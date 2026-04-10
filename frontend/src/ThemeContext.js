import React, { createContext, useContext, useState, useEffect } from 'react';

export const themes = {
  dark: {
    bg: '#08080c',
    bgSecondary: '#0f0f18',
    bgCard: 'rgba(255,255,255,0.03)',
    bgCardHover: 'rgba(255,255,255,0.05)',
    border: 'rgba(255,255,255,0.06)',
    borderStrong: 'rgba(255,255,255,0.1)',
    text: '#f0ede8',
    textSecondary: 'rgba(255,255,255,0.55)',
    textMuted: 'rgba(255,255,255,0.3)',
    textBody: 'rgba(240,237,232,0.8)',
    navBg: 'rgba(8,8,12,0.8)',
    navBgScrolled: 'rgba(8,8,12,0.96)',
    inputBg: 'rgba(255,255,255,0.04)',
    inputBorder: 'rgba(255,255,255,0.08)',
    accent: '#ff6b35',
    accentLight: 'rgba(255,107,53,0.1)',
    accentBorder: 'rgba(255,107,53,0.25)',
    isDark: true,
  },
  light: {
    bg: '#faf9f7',
    bgSecondary: '#f0ede8',
    bgCard: '#ffffff',
    bgCardHover: '#f7f5f2',
    border: 'rgba(0,0,0,0.07)',
    borderStrong: 'rgba(0,0,0,0.12)',
    text: '#1a1208',
    textSecondary: 'rgba(26,18,8,0.6)',
    textMuted: 'rgba(26,18,8,0.35)',
    textBody: 'rgba(26,18,8,0.75)',
    navBg: 'rgba(250,249,247,0.85)',
    navBgScrolled: 'rgba(250,249,247,0.97)',
    inputBg: 'rgba(0,0,0,0.03)',
    inputBorder: 'rgba(0,0,0,0.1)',
    accent: '#e85d20',
    accentLight: 'rgba(232,93,32,0.08)',
    accentBorder: 'rgba(232,93,32,0.2)',
    isDark: false,
  },
};

const ThemeContext = createContext({ theme: themes.dark, mode: 'dark', toggle: () => {} });

export function ThemeProvider({ children }) {
  const [mode, setMode] = useState(() => localStorage.getItem('theme') || 'dark');
  const theme = themes[mode];

  const toggle = () => {
    const next = mode === 'dark' ? 'light' : 'dark';
    setMode(next);
    localStorage.setItem('theme', next);
  };

  // Apply to body background
  useEffect(() => {
    document.body.style.background = theme.bg;
    document.body.style.color = theme.text;
  }, [theme]);

  return (
      <ThemeContext.Provider value={{ theme, mode, toggle }}>
        {children}
      </ThemeContext.Provider>
  );
}

export const useTheme = () => useContext(ThemeContext);