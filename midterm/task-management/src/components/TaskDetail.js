import React from 'react';
import 'bootstrap/dist/css/bootstrap.min.css';

const TaskDetail = ({ task, onClose }) => {
  if (!task) {
    return null;
  }

  return (
    <div className="modal show" style={{ display: 'block', backgroundColor: 'rgba(0, 0, 0, 0.5)' }}>
      <div className="modal-dialog">
        <div className="modal-content">
          <div className="modal-header">
            <h5 className="modal-title">Task Details</h5>
            <button type="button" className="btn-close" onClick={onClose}></button>
          </div>
          <div className="modal-body">
            <h4>{task.title}</h4>
            <p><strong>Description:</strong> {task.description}</p>
            <p><strong>Status:</strong> {task.status}</p>
            {task.Users && task.Users.length > 0 && (
              <div>
                <strong>Assigned Users:</strong>
                <ul>
                  {task.Users.map(user => (
                    <li key={user.id}>{user.name}</li>
                  ))}
                </ul>
              </div>
            )}
          </div>
          <div className="modal-footer">
            <button type="button" className="btn btn-secondary" onClick={onClose}>
              Close
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default TaskDetail;
