import React, { ChangeEvent, Fragment } from "react";

import { FormControl, FormControlLabel, FormLabel, Radio, RadioGroup } from "@mui/material";
import { useTranslation } from "react-i18next";

import { SecondFactorMethod } from "@models/Methods";
import { toMethod2FA } from "@services/UserInfo";

interface Props {
    id: string;
    methods: SecondFactorMethod[];
    method: SecondFactorMethod;
    name: string;
    handleMethodChanged: (event: ChangeEvent<HTMLInputElement>) => void;
}

const TwoFactorAuthenticationOptionsMethodsRadioGroup = function (props: Props) {
    const { t: translate } = useTranslation("settings");

    return (
        <FormControl>
            <FormLabel>{translate(props.name)}</FormLabel>
            <RadioGroup value={toMethod2FA(props.method)} onChange={props.handleMethodChanged} row>
                {props.methods.map((value, index) => {
                    const v = toMethod2FA(value);

                    switch (value) {
                        case SecondFactorMethod.WebAuthn:
                            return (
                                <FormControlLabel
                                    id={`method-${props.id}-default-webauthn`}
                                    control={<Radio />}
                                    label={translate("WebAuthn")}
                                    key={index}
                                    value={v}
                                />
                            );
                        case SecondFactorMethod.TOTP:
                            return (
                                <FormControlLabel
                                    id={`method-${props.id}-default-one-time-password`}
                                    control={<Radio />}
                                    label={translate("One-Time Password")}
                                    key={index}
                                    value={v}
                                />
                            );
                        case SecondFactorMethod.MobilePush:
                            return (
                                <FormControlLabel
                                    id={`method-${props.id}-default-duo`}
                                    control={<Radio />}
                                    label={translate("Mobile Push")}
                                    key={index}
                                    value={v}
                                />
                            );
                        default:
                            return <Fragment />;
                    }
                })}
            </RadioGroup>
        </FormControl>
    );
};

export default TwoFactorAuthenticationOptionsMethodsRadioGroup;
