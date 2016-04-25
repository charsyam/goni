// React
import React from 'react';
import { connect } from 'react-redux';

// Actions
import Actions from 'actions/Projects';

// Components
import Header from 'components/Header';

class Projects extends React.Component {
  componentDidMount() {
    const { dispatch } = this.props;
    dispatch(Actions.getProjects());
  }

  _handleProjectClick(e, project) {
    const { dispatch } = this.props;
    dispatch(Actions.enterProject(project));
  }

  _renderProjects() {
    const { projects } = this.props;
    return projects.map((project) => {
      const keyUrl = `/setup/project/${project.id}`;
      return (
        <div key={project.id} className="col-xs-12 col-sm-6 col-md-6 col-lg-4">
          <div className="project-card" onClick={(e) => this._handleProjectClick(e, project)}>
            {project.is_plus ? <p className="isplus">Goni+</p> : <p className="isplus">Goni</p>}
            <p className="title">{project.name}</p>
            <div className="tag">
              <a href={keyUrl}>APIKEY : {project.apikey}</a>
            </div>
          </div>
        </div>
      );
    });
  }

  _renderNewProjectBtn() {
    return (
      <div className="col-xs-12 col-sm-6 col-md-6 col-lg-4">
        <div className="project-card">
          <p className="isplus">NEW PROJECT</p>
          <p className="title">+ 프로젝트 추가</p>
          <div className="tag">
            <a href="https://github.com/layer123/goni">QUICKSTART GUIDE</a>
          </div>
        </div>
      </div>
    );
  }

  _renderLayout() {
    const { fetching } = this.props;
    return (
      <div className="child">
        <div className="dashboard">
          <div className="dashboard-header">
            <h1>PROJECTS { fetching ? <i className="fa fa-circle-o-notch fa-spin" aria-hidden="true" /> : null }</h1>
          </div>
          <div className="row">
            {this._renderProjects()}
            {this._renderNewProjectBtn()}
          </div>
        </div>
      </div>
    );
  }

  render() {
    return (
      <div>
        <Header page="projects" />
        {this._renderLayout()}
      </div>
    );
  }
}

Projects.propTypes = {
  currentUser: React.PropTypes.object,
  dispatch: React.PropTypes.func.isRequired,
  fetching: React.PropTypes.bool.isRequired,
  projects: React.PropTypes.array,
};

const mapStateToProps = (state) => ({
  currentUser: state.auth.currentUser,
  projects: state.projects.list,
  fetching: state.projects.fetching,
});

export default connect(mapStateToProps)(Projects);
