import React from 'react';
import { connect } from 'react-redux';
import background from '../../assets/Background.svg';
import Header from '../../components/header/Header'
import './login.css';
import LoginForm from './LoginForm';
import { set, get } from "automate-redux";

function Login(props) {
  return (
    <div className="wrapper" >
      <div className="image">
        <img className="background" src={background} />
        <div className="rText">
          <Header name="Welcome to Space Cloud !" color="#FFF" fontSize="28px" />
          <p className="text">Develop great applications without having to write backend code.
              Focus more on business and less on technology.</p>
        </div>
      </div>
      <LoginForm formState={props.formState} isLoading={props.isLoading} updateFormState={props.updateFormState} handleSubmit={props.handleSubmit} />
    </div>
  )
}

const mapStateToProps = (state) => {
  return {
    formState: state.uiState.login.formState,
    isLoading: state.uiState.login.isLoading,
  }
}

const mapDispatchToProps = (dispatch) => {
  return {
    updateFormState: (fields) => {
      console.log(fields)
      dispatch(set("uiState.login.formState", fields))
    },
    handleSubmit: () => {
      console.log('Submitted')
    }
  }
}

export default connect(mapStateToProps, mapDispatchToProps)(Login);
