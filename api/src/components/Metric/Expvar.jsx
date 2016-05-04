// React
import React from 'react';
import { connect } from 'react-redux';

// Actions
import Actions from 'actions/Metric';

// Components
import { Empty, Error, Loading, tooltipFormat } from './Common';
import { ResponsiveLineChart } from 'core/chart';
import Select from 'react-select';

// Constants
import { METRIC_CHANGE_INSTANCE } from 'constants/metric';

class Expvar extends React.Component {
  componentDidMount() {
    const { currentProject, dispatch } = this.props;
    dispatch(Actions.getInstances(currentProject.apikey, 'expvar'));
  }

  _changeInstance(v) {
    const { currentDuration, currentInstance, currentProject } = this.props;
    const { dispatch } = this.props;
    dispatch({
      type: METRIC_CHANGE_INSTANCE,
      instance: v,
    });
    dispatch(Actions.getGoniPlus(currentProject.apikey, 'expvar',
      currentInstance, currentDuration));
  }

  _renderChart(title, unit, mod) {
    const { currentDuration, currentInstance, currentProject } = this.props;
    const { dispatch, errored, fetchedData, fetching } = this.props;
    if (!currentInstance) {
      return (
        <Error title={title} msg="인스턴스를 선택해주세요" />
      );
    }
    if (!fetchedData) {
      if (!fetching && !errored) {
        dispatch(Actions.getGoniPlus(currentProject.apikey, 'expvar',
          currentInstance, currentDuration));
      }
      return (
        <Loading title={title} fetching={fetching} />
      );
    }
    if (!(title in fetchedData)) {
      return (
        <Empty title={title} />
      );
    }
    const dataLen = fetchedData[title].length;
    if (dataLen === 0) {
      return (
        <Empty title={title} />
      );
    }
    const chartData = [];
    const parsed = {
      name: title,
      values: [],
    };
    for (let i = 0; i < dataLen; i++) {
      parsed.values.push({
        x: new Date(fetchedData[title][i].time),
        y: fetchedData[title][i][title] / mod,
      });
    }
    chartData.push(parsed);
    return (
      <div>
        <div className="chart-wrapper-header">{title}</div>
        <ResponsiveLineChart data={chartData} tooltipFormat={(v) => tooltipFormat(v, unit)} />
      </div>
    );
  }

  render() {
    const { currentInstance, fetchedInstances, instanceFetching } = this.props;
    return (
      <div>
        <Select name="instance" options={fetchedInstances} onChange={::this._changeInstance}
          isLoading={instanceFetching} value={currentInstance} placeholder="Instance" />
        {this._renderChart('alloc', 'KB', 1024)}
        {this._renderChart('heapalloc', 'KB', 1024)}
        {this._renderChart('heapinuse', 'KB', 1024)}
        {this._renderChart('numgc', '', 1)}
        {this._renderChart('pausetotalns', 'ms', 1000000)}
        {this._renderChart('sys', 'KB', 1024)}
      </div>
    );
  }
}

Expvar.propTypes = {
  currentDuration: React.PropTypes.string,
  currentInstance: React.PropTypes.string,
  currentProject: React.PropTypes.object,
  dispatch: React.PropTypes.func.isRequired,
  errored: React.PropTypes.bool,
  fetchedData: React.PropTypes.object,
  fetchedInstances: React.PropTypes.Array,
  fetching: React.PropTypes.bool,
  instanceFetching: React.PropTypes.bool,
};

const mapStateToProps = (state) => ({
  currentDuration: state.project.currentDuration,
  currentInstance: state.metric.currentInstance,
  currentProject: state.project.currentProject,
  errored: state.metric.errored,
  fetchedData: state.metric.fetchedData,
  fetchedInstances: state.metric.fetchedInstances,
  fetching: state.metric.fetching,
  instanceFetching: state.metric.instanceFetching,
});

export default connect(mapStateToProps)(Expvar);
