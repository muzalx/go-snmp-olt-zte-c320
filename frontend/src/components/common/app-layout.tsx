import { NavLink, Outlet } from "react-router-dom";

const links = [
  { to: "/health", label: "Health" },
  { to: "/onu", label: "ONU Explorer" },
  { to: "/cache", label: "Cache Tools" },
];

export function AppLayout() {
  return (
    <div className="layout">
      <header className="header">
        <h1>OLT ZTE C320 Monitor</h1>
        <nav>
          {links.map((link) => (
            <NavLink
              key={link.to}
              to={link.to}
              className={({ isActive }) => `nav-link ${isActive ? "active" : ""}`}
            >
              {link.label}
            </NavLink>
          ))}
        </nav>
      </header>
      <main className="container">
        <Outlet />
      </main>
    </div>
  );
}
