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
        <div>
          <h1 className="brand-title">AMCNet Monitoring OLT</h1>
          <p className="brand-subtitle">Network Operations Dashboard</p>
        </div>
        <nav aria-label="Main navigation">
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
