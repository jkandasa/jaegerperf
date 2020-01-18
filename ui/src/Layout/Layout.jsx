import React from "react";
import { Layout, Menu, Row, Col } from "antd";
import { withRouter } from "react-router";
import { Route, Redirect, Switch } from "react-router-dom";

import "./Layout.css";

import { routes, hiddenRoutes } from "../Services/Routes";

const { Header, Content } = Layout;
const contentMargin = { xs: 0, sm: 1, md: 1, lg: 1 }; // margin left + margin right = *2
const contentBody = { xs: 24, sm: 22, md: 22, lg: 22 };

class PageLayout extends React.Component {
  //onMenuSelect = ({item, key, keyPath, domEvent}) => {
  onMenuSelect = data => {
    //console.log(item, key, keyPath, domEvent)
    //this.setState({
    //selectedMenuKey: data.key
    //});
    const { history } = this.props;
    history.push(data.item.props.item.to);
  };

  navigateTo = path => {
    const { history } = this.props;
    history.push(path);
  };

  renderContent = () => {
    const allRoutes = [];
    routes.forEach(item => {
      allRoutes.push(
        <Route key={item.to} exact path={item.to} component={item.component} />
      );
    });
    hiddenRoutes.forEach(item => {
      allRoutes.push(
        <Route key={item.to} exact path={item.to} component={item.component} />
      );
    });
    return (
      <Switch>
        {allRoutes}
        <Redirect from="*" to="/" key="default-route" />
      </Switch>
    );
  };

  render = () => {
    const { location } = this.props
    let menuSelection = ""
    routes.forEach(r => {
      if (location.pathname.startsWith(r.to)) {
        menuSelection = r.id
      }
    })
    return (
      <Layout className="layout">
        <Header>
          <div className="title">JaegerPerf</div>
          <Menu
            theme="dark"
            mode="horizontal"
            // defaultSelectedKeys={["2"]}
            selectedKeys={menuSelection}
            onClick={this.onMenuSelect}
            selectable={false}
            style={{ lineHeight: "47px" }}
          >
            {routes.map(m => {
              return (
                <Menu.Item key={m.id} item={m}>
                  {m.title}
                </Menu.Item>
              );
            })}
          </Menu>
        </Header>
        <Content>
          <Row>
            <Col {...contentMargin} />
            <Col {...contentBody}>{this.renderContent()}</Col>
            <Col {...contentMargin} />
          </Row>
        </Content>
      </Layout>
    );
  };
}

export default withRouter(PageLayout);
